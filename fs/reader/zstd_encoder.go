package reader

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/zstd"
	"pnxlr.eu.org/roll/fs/util"
)

type ZSTDEncoder struct {
	br    *BlockReader
	brbuf [util.MiB]byte
	buf   *bytes.Buffer
	enc   *zstd.Encoder
}

func NewZSTDEncoder(br *BlockReader) *ZSTDEncoder {
	enc, _ := zstd.NewWriter(nil,
		zstd.WithEncoderLevel(zstd.SpeedBetterCompression))
	return &ZSTDEncoder{
		br:  br,
		buf: &bytes.Buffer{},
		enc: enc,
	}
}

func (r *ZSTDEncoder) Read(p []byte) (int, error) {
	if !r.br.HeaderDone() {
		return r.br.Read(p)
	} else if r.buf.Len() > 0 {
		return r.buf.Read(p)
	}

	n, err := r.br.Read(r.brbuf[:])
	if err != nil && err != io.EOF {
		return 0, err
	} else if n > 0 {
		r.buf.Write(r.enc.EncodeAll(r.brbuf[:n], nil))
		n, _ := r.buf.Read(p)
		return n, err
	}

	return 0, err
}

func (r *ZSTDEncoder) Close() error {
	r.enc.Close()
	return r.br.Close()
}

package reader

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

type ZSTDDecoder struct {
	r   io.ReadCloser
	dec *zstd.Decoder
}

func NewZSTDDecoder(r io.ReadCloser) *ZSTDDecoder {
	dec, _ := zstd.NewReader(r)
	return &ZSTDDecoder{r: r, dec: dec}
}

func (r *ZSTDDecoder) Read(p []byte) (int, error) {
	return r.dec.Read(p)
}

func (r *ZSTDDecoder) Close() error {
	r.dec.Close()
	return r.r.Close()
}

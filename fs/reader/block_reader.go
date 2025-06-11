package reader

import (
	"io"

	"pnxlr.eu.org/roll/fs/util"
	"pnxlr.eu.org/roll/util/log"
)

type BlockReader struct {
	r           io.Reader
	totalSize   int
	totalRead   int
	lastRead    int
	progress    int
	headerBytes []byte
	headerDone  bool
	headerRead  int
}

func NewBlockReader(r io.Reader, totalSize int, headerBytes []byte) *BlockReader {
	return &BlockReader{
		r:           r,
		totalSize:   totalSize,
		headerBytes: headerBytes,
	}
}

func (r *BlockReader) Read(p []byte) (int, error) {
	if len(p) > util.MiB {
		p = p[:util.MiB]
	}

	if !r.headerDone {
		n := copy(p, r.headerBytes[r.headerRead:])
		if r.headerRead += n; r.headerRead >= len(r.headerBytes) {
			r.headerDone = true
		}
		return n, nil
	}

	n, err := r.r.Read(p)
	r.totalRead += n
	r.reportProgress()

	return n, err
}

func (r *BlockReader) reportProgress() {
	if r.totalSize <= 0 {
		return
	} else if prg := 100 * r.totalRead / r.totalSize; r.progress <
		prg && (r.lastRead+util.MiB < r.totalRead || prg == 100) {
		r.lastRead, r.progress = r.totalRead, prg
		log.Infof("Read: %d%%", prg)
	}
}

func (r *BlockReader) Close() error {
	if c, ok := r.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (r *BlockReader) HeaderDone() bool {
	return r.headerDone
}

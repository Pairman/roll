package header

import (
	"bytes"
	"fmt"
	"strconv"

	"pnxlr.eu.org/roll/fs/util"
)

type CompressionAlgoType int8

const (
	CompressionAlgoNone CompressionAlgoType = iota
	CompressionAlgoZSTD
)

func (a CompressionAlgoType) String() string {
	switch a {
	case CompressionAlgoNone:
		return "None"
	case CompressionAlgoZSTD:
		return "ZSTD"
	default:
		return strconv.Itoa(int(a))
	}
}

type CompressionSect struct { // File compression
	Algo     CompressionAlgoType // Algorithm
	InfoSize int16               // Information length
	Info     []byte              // Information
}

func (s CompressionSect) String() string {
	return fmt.Sprintf("CompressionSect{Algo: %v, InfoSize: %v, Info: %v}",
		s.Algo, s.InfoSize, s.Info)
}

func NewCompressionSect(algo CompressionAlgoType) (*CompressionSect, error) {
	switch algo {
	case CompressionAlgoNone:
	case CompressionAlgoZSTD:
	default:
		return nil, fmt.Errorf("unknown compression algorithm: %v", algo)
	}
	return &CompressionSect{Algo: algo}, nil
}

func (s *CompressionSect) ToBytes() []byte {
	buf := bytes.Buffer{}
	buf.Write(util.LiteralToBytes(s.Algo))
	buf.Write(util.LiteralToBytes(s.InfoSize))
	buf.Write(s.Info)
	return buf.Bytes()
}

func (s *CompressionSect) FromBytes(p []byte) error {
	s.Algo = util.LiteralFromBytes[CompressionAlgoType](p[:1])
	s.InfoSize = util.LiteralFromBytes[int16](p[1:3])
	if sectLen := s.Len(); len(p) != sectLen {
		return fmt.Errorf("info length mismatch: %d, %d", sectLen, len(p))
	}
	buff := bytes.Buffer{}
	buff.Write(p[3:s.Len()])
	s.Info = buff.Bytes()
	return nil
}

func (s *CompressionSect) Len() int {
	return 3 + int(s.InfoSize)
}

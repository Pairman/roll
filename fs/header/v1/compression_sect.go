package v1

import (
	"bytes"
	"fmt"

	"pnxlr.eu.org/roll/fs/header"
	"pnxlr.eu.org/roll/fs/util"
)

type CompressionSect struct { // File compression
	Algo     header.CompressionAlgoType // Algorithm
	InfoSize int16                      // Information length
	Info     []byte                     // Information
}

func (s CompressionSect) String() string {
	return fmt.Sprintf("CompressionSect{Algo: %v, InfoSize: %v, Info: %v}",
		s.Algo, s.InfoSize, s.Info)
}

func NewCompressionSect(
	algo header.CompressionAlgoType) (*CompressionSect, error) {
	switch algo {
	case header.CompressionAlgoNone:
	case header.CompressionAlgoZSTD:
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
	s.Algo = util.LiteralFromBytes[header.CompressionAlgoType](
		p[:header.SizeCompressionAlgoType])
	offs := header.SizeCompressionAlgoType
	s.InfoSize = util.LiteralFromBytes[int16](p[offs : offs+header.SizeInt16])
	offs += header.SizeInt16
	if sectLen := s.Len(); len(p) != sectLen {
		return fmt.Errorf("info length mismatch: %d, %d", sectLen, len(p))
	}
	buff := bytes.Buffer{}
	buff.Write(p[offs:s.Len()])
	s.Info = buff.Bytes()
	return nil
}

func (s *CompressionSect) Len() int {
	return header.SizeCompressionAlgoType + header.SizeInt16 + int(s.InfoSize)
}

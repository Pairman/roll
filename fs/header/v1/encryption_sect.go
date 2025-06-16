package v1

import (
	"bytes"
	"fmt"

	"pnxlr.eu.org/roll/fs/header"
	"pnxlr.eu.org/roll/fs/util"
)

type EncryptionSect struct { // File encryption
	Algo     header.EncryptionAlgoType // Algorithm
	InfoSize int16                     // Information length
	Info     []byte                    // Information
}

func (s EncryptionSect) String() string {
	return fmt.Sprintf("EncryptionSect{Algo: %v, InfoSize: %v, Info: %v}",
		s.Algo, s.InfoSize, s.Info)
}

func NewEncryptionSect(
	algo header.EncryptionAlgoType) (*EncryptionSect, error) {
	var size int16
	info := []byte{}
	switch algo {
	case header.EncryptionAlgoNone:
	case header.EncryptionAlgoAES256GCM:
	default:
		return nil, fmt.Errorf("unknown encryption algorithm: %v", algo)
	}
	return &EncryptionSect{Algo: algo, InfoSize: size, Info: info}, nil
}

func (s *EncryptionSect) ToBytes() []byte {
	buf := bytes.Buffer{}
	buf.Write(util.LiteralToBytes(s.Algo))
	buf.Write(util.LiteralToBytes(s.InfoSize))
	buf.Write(s.Info)
	return buf.Bytes()
}

func (s *EncryptionSect) FromBytes(p []byte) error {
	s.Algo = util.LiteralFromBytes[header.EncryptionAlgoType](
		p[:header.SizeEncryptionAlgoType])
	offs := header.SizeEncryptionAlgoType
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

func (s *EncryptionSect) Len() int {
	return header.SizeEncryptionAlgoType + header.SizeInt16 + int(s.InfoSize)
}

package header

import (
	"bytes"
	"fmt"
	"strconv"

	"pnxlr.eu.org/roll/fs/util"
)

type EncryptionAlgoType int8

const (
	EncryptionAlgoNone EncryptionAlgoType = iota
	EncryptionAlgoAES256GCM
)

func (a EncryptionAlgoType) String() string {
	switch a {
	case EncryptionAlgoNone:
		return "None"
	case EncryptionAlgoAES256GCM:
		return "AES-256-GCM"
	default:
		return strconv.Itoa(int(a))
	}
}

type EncryptionSect struct { // File encryption
	Algo     EncryptionAlgoType // Algorithm
	InfoSize int16              // Information length
	Info     []byte             // Information
}

func (s EncryptionSect) String() string {
	return fmt.Sprintf("EncryptionSect{Algo: %v, InfoSize: %v, Info: %v}",
		s.Algo, s.InfoSize, s.Info)
}

func NewEncryptionSect(algo EncryptionAlgoType) (*EncryptionSect, error) {
	var size int16
	info := []byte{}
	switch algo {
	case EncryptionAlgoNone:
	case EncryptionAlgoAES256GCM:
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
	s.Algo = util.LiteralFromBytes[EncryptionAlgoType](p[:1])
	s.InfoSize = util.LiteralFromBytes[int16](p[1:3])
	if sectLen := s.Len(); len(p) != sectLen {
		return fmt.Errorf("info length mismatch: %d, %d", sectLen, len(p))
	}
	buff := bytes.Buffer{}
	buff.Write(p[3:s.Len()])
	s.Info = buff.Bytes()
	return nil
}

func (s *EncryptionSect) Len() int {
	return 3 + int(s.InfoSize)
}

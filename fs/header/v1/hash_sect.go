package v1

import (
	"bytes"
	"fmt"
	"io"

	"github.com/zeebo/xxh3"
	"pnxlr.eu.org/roll/fs/header"
	"pnxlr.eu.org/roll/fs/util"
)

type HashSect struct { // File hashing
	Algo     header.HashAlgoType // Algorithm
	HashSize int16               // Hash length
	Hash     []byte              // Hash
}

func (s HashSect) String() string {
	return fmt.Sprintf("HashSect{Algo: %v, HashSize: %v, Hash: %x}",
		s.Algo, s.HashSize, s.Hash)
}

func NewHashSect(r io.Reader, algo header.HashAlgoType) (*HashSect, error) {
	trySeekStart(r)
	hash, err := funcHash(r, algo)
	trySeekStart(r)
	if err != nil {
		return nil, err
	}
	return &HashSect{
		Algo: algo, HashSize: int16(len(hash)), Hash: hash,
	}, nil
}

func (s *HashSect) ToBytes() []byte {
	buf := bytes.Buffer{}
	buf.Write(util.LiteralToBytes(s.Algo))
	buf.Write(util.LiteralToBytes(s.HashSize))
	buf.Write(s.Hash)
	return buf.Bytes()
}

func (s *HashSect) FromBytes(p []byte) error {
	s.Algo = util.LiteralFromBytes[header.HashAlgoType](
		p[:header.SizeHashAlgoType])
	offs := header.SizeHashAlgoType
	s.HashSize = util.LiteralFromBytes[int16](p[offs : offs+header.SizeInt16])
	offs += header.SizeInt16
	if sectLen := s.Len(); len(p) != sectLen {
		return fmt.Errorf("buffer length mismatch: %d, %d", sectLen, len(p))
	}
	buff := bytes.Buffer{}
	buff.Write(p[offs:s.Len()])
	s.Hash = buff.Bytes()
	return nil
}

func (s *HashSect) Len() int {
	return header.SizeHashAlgoType + header.SizeInt16 + int(s.HashSize)
}

func (s *HashSect) Verify(r io.Reader) (bool, []byte, error) {
	trySeekStart(r)
	hash, err := funcHash(r, s.Algo)
	trySeekStart(r)
	if err != nil {
		return false, nil, err
	}
	return bytes.Equal(hash, s.Hash), hash, nil
}

func funcHash(r io.Reader, algo header.HashAlgoType) ([]byte, error) {
	hash := []byte{}
	switch algo {
	case header.HashAlgoNone:
	case header.HashAlgoXXH3:
		hasher := xxh3.New()
		io.Copy(hasher, r) // Calls hasher.Write(), which never errors
		hash = util.LiteralToBytes(hasher.Sum64())
	default:
		return nil, fmt.Errorf("unknown hash algorithm: %v", algo)
	}
	return hash, nil
}

func trySeekStart(r any) {
	if s, ok := r.(io.Seeker); ok {
		s.Seek(0, io.SeekStart)
	}
}

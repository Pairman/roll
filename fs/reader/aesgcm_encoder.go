package reader

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io"

	"pnxlr.eu.org/roll/fs/util"
)

const gcmNonceSize = 12

type AESGCMEncoder struct {
	br    *BlockReader
	brbuf [util.MiB]byte
	buf   *bytes.Buffer
	aead  cipher.AEAD
	iv    [gcmNonceSize]byte
	cnt   uint32
}

func NewAESGCMEncoder(br *BlockReader, key, iv []byte) *AESGCMEncoder {
	if len(key) != 32 {
		panic("Key must be 32 bytes")
	} else if len(iv) != gcmNonceSize {
		panic("IV must be 12 bytes")
	}
	block, _ := aes.NewCipher(key)
	aead, _ := cipher.NewGCM(block)
	r := &AESGCMEncoder{
		br:   br,
		buf:  &bytes.Buffer{},
		aead: aead,
		cnt:  binary.BigEndian.Uint32(iv[gcmNonceSize-4:]),
	}
	copy(r.iv[:], iv)
	return r
}

func (r *AESGCMEncoder) Read(p []byte) (int, error) {
	if !r.br.HeaderDone() {
		return r.br.Read(p)
	} else if r.buf.Len() > 0 {
		return r.buf.Read(p)
	}

	n, err := r.br.Read(r.brbuf[:])
	if err != nil && err != io.EOF {
		return 0, err
	} else if n > 0 {
		var nonce [gcmNonceSize]byte
		copy(nonce[:], r.iv[:])
		binary.BigEndian.PutUint32(nonce[gcmNonceSize-4:], r.cnt)
		r.cnt += 1

		var ctLen [4]byte
		ct := r.aead.Seal(nil, nonce[:], r.brbuf[:n], nil)
		binary.BigEndian.PutUint32(ctLen[:], uint32(len(ct)))
		r.buf.Write(ctLen[:])

		r.buf.Write(ct)

		if err == io.EOF {
			err = nil
		}
		m, _ := r.buf.Read(p)
		return m, err
	}

	return 0, io.EOF
}

func (r *AESGCMEncoder) Close() error {
	return r.br.Close()
}

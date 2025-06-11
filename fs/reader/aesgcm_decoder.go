package reader

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io"

	"pnxlr.eu.org/roll/fs/util"
)

type AESGCMDecoder struct {
	r     io.ReadCloser
	rbuf  [util.MiB + 16]byte
	buf   *bytes.Buffer
	aead  cipher.AEAD
	nonce [gcmNonceSize]byte
	cnt   uint32
}

func NewAESGCMDecoder(r io.ReadCloser, key, iv []byte) *AESGCMDecoder {
	block, _ := aes.NewCipher(key)
	aead, _ := cipher.NewGCM(block)
	d := &AESGCMDecoder{
		r: r, buf: &bytes.Buffer{}, aead: aead,
		cnt: binary.BigEndian.Uint32(iv[aead.NonceSize()-4:]),
	}
	copy(d.nonce[:], iv)
	return d
}

func (r *AESGCMDecoder) Read(p []byte) (int, error) {
	if r.buf.Len() > 0 {
		return r.buf.Read(p)
	}

	var ctLenBuf [4]byte
	if _, err := io.ReadFull(r.r, ctLenBuf[:]); err != nil {
		return 0, err
	}
	ctLen := binary.BigEndian.Uint32(ctLenBuf[:])
	if _, err := io.ReadFull(r.r, r.rbuf[:ctLen]); err != nil {
		return 0, err
	}

	binary.BigEndian.PutUint32(r.nonce[gcmNonceSize-4:], r.cnt)
	r.cnt++
	pt, err := r.aead.Open(nil, r.nonce[:], r.rbuf[:ctLen], nil)
	if err != nil {
		return 0, err
	}
	r.buf.Write(pt)
	return r.buf.Read(p)
}

func (r *AESGCMDecoder) Close() error {
	return r.r.Close()
}

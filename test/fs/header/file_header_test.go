package header_test

import (
	"bytes"
	"os"
	"testing"

	"pnxlr.eu.org/roll/fs/header"
)

func TestFileHeader(t *testing.T) {
	f, _ := os.CreateTemp("", "测试.txt")
	defer os.Remove(f.Name())
	defer f.Close()
	f.Write([]byte("测试 test"))

	fs, err := header.NewFileSect(f)
	if err != nil {
		t.Errorf("NewFileSect() error: %v", err)
	}

	hs, err := header.NewHashSect(f, 1)
	if err != nil {
		t.Errorf("NewHashSect() error: %v", err)
	}

	cs, err := header.NewCompressionSect(header.CompressionAlgoZSTD)
	if err != nil {
		t.Errorf("NewCompressionSect() error: %v", err)
	}

	es, err := header.NewEncryptionSect(header.EncryptionAlgoAES256GCM)
	if err != nil {
		t.Errorf("NewEncryptionSect() error: %v", err)
	}

	fh := header.NewFileHeader(fs, hs, cs, es)
	if fh.Len() != len(fh.ToBytes()) {
		t.Errorf("Len() error: %v, %v", fh.Len(), len(fh.ToBytes()))
	}

	nfh := &header.FileHeader{}
	if err := nfh.FromBytes(fh.ToBytes()); err != nil {
		t.Errorf("FromBytes() error: %v", err)
	}
	if !bytes.Equal(nfh.ToBytes(), fh.ToBytes()) {
		t.Errorf("ToBytes() error: %v", false)
	}
}

package header_test

import (
	"bytes"
	"os"
	"testing"

	header "pnxlr.eu.org/roll/fs/header"
	headerV1 "pnxlr.eu.org/roll/fs/header/v1"
)

func TestFileHeader(t *testing.T) {
	f, _ := os.CreateTemp("", "测试.txt")
	defer os.Remove(f.Name())
	defer f.Close()
	f.Write([]byte("测试 test"))

	fs, err := headerV1.NewFileSect(f)
	if err != nil {
		t.Errorf("NewFileSect() error: %v", err)
	}

	hs, err := headerV1.NewHashSect(f, 1)
	if err != nil {
		t.Errorf("NewHashSect() error: %v", err)
	}

	cs, err := headerV1.NewCompressionSect(header.CompressionAlgoZSTD)
	if err != nil {
		t.Errorf("NewCompressionSect() error: %v", err)
	}

	es, err := headerV1.NewEncryptionSect(header.EncryptionAlgoAES256GCM)
	if err != nil {
		t.Errorf("NewEncryptionSect() error: %v", err)
	}

	fh := headerV1.NewFileHeader(fs, hs, cs, es)
	if fh.Len() != len(fh.ToBytes()) {
		t.Errorf("Len() error: %v, %v", fh.Len(), len(fh.ToBytes()))
	}

	nfh := &headerV1.FileHeader{}
	if err := nfh.FromBytes(fh.ToBytes()); err != nil {
		t.Errorf("FromBytes() error: %v", err)
	}
	if !bytes.Equal(nfh.ToBytes(), fh.ToBytes()) {
		t.Errorf("ToBytes() error: %v", false)
	}
}

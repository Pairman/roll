package reader_test

import (
	"bytes"
	"crypto/rand"
	"os"
	"testing"

	headerV1 "pnxlr.eu.org/roll/fs/header/v1"
	"pnxlr.eu.org/roll/fs/reader"
	"pnxlr.eu.org/roll/fs/util"
)

func TestAESGCMEncoderDecoder(t *testing.T) {
	const TEXT = "cafebabe"
	f, _ := os.CreateTemp("", "TestAESGCMEncoderDecoder.txt")
	defer os.Remove(f.Name())
	defer f.Close()
	f.Write(bytes.Repeat([]byte(TEXT), util.MiB/len(TEXT)/2*3))

	fh := headerV1.NewFileHeaderFromFile(f)
	br := reader.NewBlockReader(f, 0, fh.ToBytes())
	defer br.Close()

	key, iv := make([]byte, 32), make([]byte, 12)
	rand.Read(key)
	rand.Read(iv)
	ar := reader.NewAESGCMEncoder(br, key, iv)
	defer ar.Close()

	p := make([]byte, util.MiB)
	n, err := ar.Read(p)
	if err != nil {
		t.Errorf("Read() error: %v, %v", n, err)
	} else if !bytes.Equal(fh.ToBytes(), p[:n]) {
		t.Errorf("Read() error: expected %v, got %v", fh.ToBytes(), p[:n])
	}

	ad := reader.NewAESGCMDecoder(ar, key, iv)
	defer ad.Close()
	buf := &bytes.Buffer{}
	if n, err := buf.ReadFrom(ad); err != nil {
		t.Errorf("ReadFrom() error: %v, %v", n, err)
	}
	p = buf.Bytes()

	if match, _, err := fh.HashSect.Verify(
		bytes.NewReader(p)); match != true || err != nil {
		t.Errorf("Verify() error: %v, %v", match, err)
	}
}

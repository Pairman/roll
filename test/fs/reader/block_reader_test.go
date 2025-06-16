package reader_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	headerV1 "pnxlr.eu.org/roll/fs/header/v1"
	"pnxlr.eu.org/roll/fs/reader"
	"pnxlr.eu.org/roll/fs/util"
)

func TestBlockReader(t *testing.T) {
	const TEXT = "cafebabe"
	f, _ := os.CreateTemp("", "TestBlockReader.txt")
	defer os.Remove(f.Name())
	defer f.Close()
	f.Write(bytes.Repeat([]byte("cafebabe"), util.MiB/len(TEXT)/2*3))

	fh := headerV1.NewFileHeaderFromFile(f)
	br := reader.NewBlockReader(f, int(fh.FileSect.FileSize), fh.ToBytes())
	defer br.Close()

	p := make([]byte, util.MiB)
	n, err := br.Read(p)
	if err != nil {
		t.Errorf("Read() error: %v, %v", n, err)
	}
	if !bytes.Equal(fh.ToBytes(), p[:n]) {
		t.Errorf("Read() error: expected %v, got %v", fh.ToBytes(), p[:n])
	}

	p, err = io.ReadAll(br)
	if err != nil {
		t.Errorf("ReadAll() error: %v, %v", len(p), err)
	}
	if s, _ := f.Stat(); len(p) != int(s.Size()) {
		t.Errorf("ReadAll() error: length expected %v, got %v", len(p), s.Size())
	}

	if match, _, err := fh.HashSect.Verify(
		bytes.NewReader(p)); match != true || err != nil {
		t.Errorf("Verify() error: %v, %v", match, err)
	}
}

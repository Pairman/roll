package reader_test

import (
	"bytes"
	"os"
	"testing"

	headerV1 "pnxlr.eu.org/roll/fs/header/v1"
	"pnxlr.eu.org/roll/fs/reader"
	"pnxlr.eu.org/roll/fs/util"
)

func TestZSTDEncoderDecoder(t *testing.T) {
	const TEXT = "cafebabe"
	f, _ := os.CreateTemp("", "TestZSTDEncoderDecoder.txt")
	defer os.Remove(f.Name())
	defer f.Close()
	f.Write(bytes.Repeat([]byte(TEXT), util.MiB/len(TEXT)/2*3))

	fh := headerV1.NewFileHeaderFromFile(f)
	br := reader.NewBlockReader(f, 0, fh.ToBytes())
	defer br.Close()

	zr := reader.NewZSTDEncoder(br)
	defer zr.Close()

	p := make([]byte, util.MiB)
	n, err := zr.Read(p)
	if err != nil {
		t.Errorf("Read() error: %v, %v", n, err)
	} else if !bytes.Equal(fh.ToBytes(), p[:n]) {
		t.Errorf("Read() error: expected %v, got %v", fh.ToBytes(), p[:n])
	}

	zd := reader.NewZSTDDecoder(zr)
	defer zd.Close()
	buf := &bytes.Buffer{}
	if n, err := buf.ReadFrom(zd); err != nil {
		t.Errorf("ReadFrom() error: %v, %v", n, err)
	}
	p = buf.Bytes()

	if match, _, err := fh.HashSect.Verify(
		bytes.NewReader(p)); match != true || err != nil {
		t.Errorf("Verify() error: %v, %v", match, err)
	}
}

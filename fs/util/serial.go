package util

import (
	"bytes"
	"encoding/binary"
)

type LiteralTypes interface {
	~int8 | ~int16 | ~int32 | ~int64 |
		~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~complex64 | ~complex128
}

func LiteralToBytes[T LiteralTypes](v T) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func LiteralFromBytes[T LiteralTypes](b []byte) T {
	var val T
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &val)
	if err != nil {
		panic(err)
	}
	return val
}

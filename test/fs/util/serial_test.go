package util_test

import (
	"testing"

	"pnxlr.eu.org/roll/fs/util"
)

func TestSzjcUploader(t *testing.T) {
	var buf []byte
	var i64 int64 = -694206942069420
	buf = util.LiteralToBytes(i64)
	ni64 := util.LiteralFromBytes[int64](buf)
	if i64 != ni64 {
		t.Errorf("Serialization error for %T: %v != %v", i64, i64, ni64)
	}

	var u16 uint16 = 42069
	buf = util.LiteralToBytes(u16)
	nu16 := util.LiteralFromBytes[uint16](buf)
	if u16 != nu16 {
		t.Errorf("Serialization error for %T: %v != %v", u16, u16, nu16)
	}

	var f64 float64 = 420.69
	buf = util.LiteralToBytes(f64)
	nf64 := util.LiteralFromBytes[float64](buf)
	if f64 != nf64 {
		t.Errorf("Serialization error for %T: %v != %v", f64, f64, nf64)
	}

	var c64 complex64 = 420.69 + 69.42i
	buf = util.LiteralToBytes(c64)
	nc64 := util.LiteralFromBytes[complex64](buf)
	if c64 != nc64 {
		t.Errorf("Serialization error for %T: %v != %v", c64, c64, nc64)
	}
}

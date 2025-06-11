package util_test

import (
	"testing"

	"pnxlr.eu.org/roll/net/util"
)

func TestObjectIDToStatus(t *testing.T) {
	if _, err := util.ObjectIDToStatus(
		"c1bb079b7bcb268e7a2817bd5263b4c5"); err != nil {
		t.Errorf("ObjectIDToStatus() error: %v", err)
	}
}

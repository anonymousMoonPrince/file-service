package utils

import (
	"testing"
)

func TestSplit(t *testing.T) {
	parts := Split(333, 4)
	if len(parts) != 4 {
		t.Fatal("incorrect part count")
	}

	if parts[0] != 83 || parts[1] != 83 || parts[2] != 83 || parts[3] != 84 {
		t.Fatal("incorrect part values")
	}
}

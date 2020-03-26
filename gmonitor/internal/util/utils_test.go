package util

import (
	"testing"
)

func TestIsLegal(t *testing.T) {
	if isLegal("123") != true {
		t.Error("Numbers")
	}
	if isLegal("abc") != true {
		t.Error("Lower case")
	}
	if isLegal("ABC") != true {
		t.Error("Upper case")
	}
	if isLegal("ABC_123abc_") != true {
		t.Error("_____")
	}
	if isLegal("!@#") != false {
		t.Error("Include illegal letter")
	}
	if isLegal("testx_23-123") != false {
		t.Error("Include illegal letter")
	}
}

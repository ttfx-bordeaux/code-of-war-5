package util

import (
	"os"
	"testing"
)

func TestLoadArg(t *testing.T) {
	os.Args = []string{"-test", "plop"}
	if v := LoadArg("-test", ""); v != "plop" {
		t.Fail()
	}
}

func TestLoadArgWithDefaultValue(t *testing.T) {
	if v := LoadArg("-fail", "defaultValue"); v != "defaultValue" {
		t.Fail()
	}
}

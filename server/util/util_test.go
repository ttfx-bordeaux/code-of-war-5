package util

import (
	"os"
	"testing"
)

func TestLoadArg(t *testing.T) {
	os.Args = []string{"-t", "--test", "plop"}
	if v := LoadArg("--test", "", "autre"); v != "plop" {
		t.Fail()
	}
}

func TestLoadArgWithDefaultValue(t *testing.T) {
	if v := LoadArg("-f", "--fail", "defaultValue"); v != "defaultValue" {
		t.Fail()
	}
}

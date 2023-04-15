package utils

import (
	"os"
	"testing"
)

func TestUtils_IsPathValid(t *testing.T) {
	file, err := os.OpenFile("test", os.O_CREATE, 0755)
	if err != nil {
		t.Fatal(err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	isPathValid := IsPathValid("test")
	if !isPathValid {
		t.Error("file has invalid path")
	}

	err = os.Remove("test")
	if err != nil {
		t.Fatal(err)
	}
}

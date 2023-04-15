package writer

import (
	"os"
	"testing"
)

func TestWriter_ToFile(t *testing.T) {
	_, err := ToFile("test", []byte("Hello, World!"))
	if err != nil {
		t.Error(err)
	}

	err = os.Remove("test")
	if err != nil {
		t.Error(err)
	}
}

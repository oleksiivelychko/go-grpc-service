package request

import (
	"testing"
)

func TestRequest_GET(t *testing.T) {
	bytesArr, err := GET("https://go.dev")
	if err != nil {
		t.Fatal(err)
	}

	if len(bytesArr) == 0 {
		t.Fatal("unable to read data")
	}
}

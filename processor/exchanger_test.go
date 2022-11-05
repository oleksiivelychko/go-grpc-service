package processor

import (
	"fmt"
	"github.com/oleksiivelychko/go-grpc-protobuf/data"
	"testing"
)

func TestNewExchanger(t *testing.T) {
	extractor := data.NewExtractorXml(data.SourceLocal)
	e, err := NewExchanger(extractor)

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("Rates: %#v\n", e.rates)
}

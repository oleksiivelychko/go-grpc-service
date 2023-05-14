package exchanger

import (
	"github.com/oleksiivelychko/go-grpc-service/extractor"
	"testing"
)

func TestExchanger_NewProcessor(t *testing.T) {
	exchanger, err := NewProcessor(extractor.New(extractor.SourceLocal, "./../rates.xml"))
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(exchanger.rates) == 0 {
		t.Fatal("unable to process rates")
	}

	t.Logf("Rates: %#v\n", exchanger.rates)
}

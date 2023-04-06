package processor

import (
	"github.com/oleksiivelychko/go-grpc-service/extractor"
	"testing"
)

func TestProcessor_NewExchanger(t *testing.T) {
	pullerXML := extractor.NewPullerXML(extractor.SourceLocal, "./../rates.xml")
	exchanger, err := NewExchanger(pullerXML)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(exchanger.rates) == 0 {
		t.Fatal("unable to process rates")
	}

	t.Logf("Rates: %#v\n", exchanger.rates)
}

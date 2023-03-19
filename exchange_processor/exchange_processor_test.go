package exchange_processor

import (
	"fmt"
	"github.com/oleksiivelychko/go-grpc-service/extractor_xml"
	"testing"
)

func TestExchangeProcessor_NewExchangeProcessor(t *testing.T) {
	extractorXML := extractor_xml.NewExtractorXML(extractor_xml.SourceLocal, "./../rates.xml")
	exchangeProcessor, err := NewExchangeProcessor(extractorXML)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(exchangeProcessor.rates) == 0 {
		t.Fatal("unable to process rates")
	}

	fmt.Printf("Rates: %#v\n", exchangeProcessor.rates)
}

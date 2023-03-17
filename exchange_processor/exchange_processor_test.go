package exchange_processor

import (
	"fmt"
	"github.com/oleksiivelychko/go-grpc-service/extractor_xml"
	"testing"
)

func TestExchangeProcessor_NewExchangeProcessor(t *testing.T) {
	extractorXML := extractor_xml.NewExtractorXML(extractor_xml.SourceLocal)
	exchangeProcessor, err := NewExchangeProcessor(extractorXML)

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("Rates: %#v\n", exchangeProcessor.rates)

	err = extractorXML.RemoveFile()
	if err != nil {
		t.Errorf(err.Error())
	}
}

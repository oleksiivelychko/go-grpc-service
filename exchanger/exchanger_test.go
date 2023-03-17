package exchanger

import (
	"fmt"
	extractor "github.com/oleksiivelychko/go-grpc-service/xml_extractor"
	"testing"
)

func TestNewExchanger(t *testing.T) {
	xmlExtractor := extractor.NewXmlExtractor(extractor.SourceLocal)
	e, err := NewExchanger(xmlExtractor)

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("Rates: %#v\n", e.rates)

	err = xmlExtractor.RemoveFile()
	if err != nil {
		t.Errorf(err.Error())
	}
}

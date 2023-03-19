package extractor_xml

import (
	"fmt"
	"github.com/oleksiivelychko/go-utils/file"
	"strconv"
	"testing"
)

const localXML = "rates.xml"

func TestExtractorXML_FetchDataFromLocalFirstTime(t *testing.T) {
	extractorXML := NewExtractorXML(SourceLocal, localXML)

	if file.DoesFileExist(localXML) {
		err := file.DeleteFile(localXML)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	err := extractorXML.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if extractorXML.source != SourceURL {
		t.Fatal("source doesn't equal to SourceURL")
	}

	testData(extractorXML, t)
}

func TestExtractorXML_FetchDataFromLocal(t *testing.T) {
	if !file.DoesFileExist(localXML) {
		t.Fatalf("file %s doesn't exist", localXML)
	}

	extractorXML := NewExtractorXML(SourceLocal, localXML)

	err := extractorXML.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if extractorXML.source != SourceLocal {
		t.Fatal("source doesn't equal to SourceURL")
	}

	testData(extractorXML, t)

	err = file.DeleteFile(localXML)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestExtractorXML_FetchDataFromURL(t *testing.T) {
	extractorXML := NewExtractorXML(SourceURL, localXML)

	err := extractorXML.FetchData()
	if err != nil {
		t.Fatal(err.Error())
	}

	if extractorXML.source != SourceURL {
		t.Fatal("source doesn't equal to SourceURL")
	}

	testData(extractorXML, t)
}

func testData(extractorXML *ExtractorXML, t *testing.T) {
	if extractorXML.RootNode.Data.Time == "" {
		t.Fatal("attribute 'time' couldn't be extracted from 'Cube' element")
	}

	for _, cube := range extractorXML.RootNode.Data.Cubes {
		rate, parseFloatErr := strconv.ParseFloat(cube.Rate, 64)
		if parseFloatErr != nil {
			t.Error(parseFloatErr)
		}
		fmt.Printf("currency: %s, rate: %f\n", cube.Currency, rate)
	}
}

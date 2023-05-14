package extractor

import (
	"errors"
	"os"
	"strconv"
	"testing"
)

const localXML = "rates.xml"

func TestExtractor_TryToExtractDataFromLocalXML(t *testing.T) {
	extractor := New(SourceLocal, localXML)

	if _, err := os.Stat(localXML); !errors.Is(err, os.ErrNotExist) {
		err = os.Remove(localXML)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	err := extractor.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if extractor.source != SourceURL {
		t.Fatal("source is not equal to SourceURL")
	}

	testData(extractor, t)
}

func TestExtractor_ExtractDataFromLocalXML(t *testing.T) {
	if _, err := os.Stat(localXML); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	extractor := New(SourceLocal, localXML)

	err := extractor.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if extractor.source != SourceLocal {
		t.Fatal("source is not equal to SourceURL")
	}

	testData(extractor, t)

	err = os.Remove(localXML)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestExtractor_ExtractDataFromURL(t *testing.T) {
	extractor := New(SourceURL, localXML)

	err := extractor.FetchData()
	if err != nil {
		t.Fatal(err.Error())
	}

	if extractor.source != SourceURL {
		t.Fatal("source is not equal to SourceURL")
	}

	testData(extractor, t)
}

func testData(extractor *XML, t *testing.T) {
	if extractor.RootNode.Data.Time == "" {
		t.Fatal("attribute 'time' could not be extracted from 'Cube' element")
	}

	for _, cube := range extractor.RootNode.Data.Cubes {
		rate, parseFloatErr := strconv.ParseFloat(cube.Rate, 64)
		if parseFloatErr != nil {
			t.Error(parseFloatErr)
			continue
		}

		t.Logf("currency: %s, rate: %f\n", cube.Currency, rate)
	}
}

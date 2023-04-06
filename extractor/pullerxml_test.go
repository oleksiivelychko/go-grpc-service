package extractor

import (
	"github.com/oleksiivelychko/go-utils/system"
	"os"
	"strconv"
	"testing"
)

const localXML = "rates.xml"

func TestExtractor_TryToPullDataFromLocalXML(t *testing.T) {
	puller := NewPullerXML(SourceLocal, localXML)

	if system.IsPathValid(localXML) {
		err := os.Remove(localXML)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	err := puller.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if puller.source != SourceURL {
		t.Fatal("source does not equal to SourceURL")
	}

	testData(puller, t)
}

func TestExtractor_PullDataFromLocalXML(t *testing.T) {
	if !system.IsPathValid(localXML) {
		t.Fatalf("file %s does not exist", localXML)
	}

	puller := NewPullerXML(SourceLocal, localXML)

	err := puller.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if puller.source != SourceLocal {
		t.Fatal("source does not equal to SourceURL")
	}

	testData(puller, t)

	err = os.Remove(localXML)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestExtractor_PullDataFromURL(t *testing.T) {
	puller := NewPullerXML(SourceURL, localXML)

	err := puller.FetchData()
	if err != nil {
		t.Fatal(err.Error())
	}

	if puller.source != SourceURL {
		t.Fatal("source does not equal to SourceURL")
	}

	testData(puller, t)
}

func testData(puller *PullerXML, t *testing.T) {
	if puller.RootNode.Data.Time == "" {
		t.Fatal("attribute 'time' could not be extracted from 'Cube' element")
	}

	for _, cube := range puller.RootNode.Data.Cubes {
		rate, parseFloatErr := strconv.ParseFloat(cube.Rate, 64)
		if parseFloatErr != nil {
			t.Error(parseFloatErr)
			continue
		}

		t.Logf("currency: %s, rate: %f\n", cube.Currency, rate)
	}
}

package data

import (
	"strconv"
	"testing"
)

func TestFetchRatesFromRemote(t *testing.T) {
	extractor := NewExtractorXml(sourceRemote)

	err := extractor.FetchRates()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if extractor.source != sourceRemote {
		t.Fatal("extracted from non-remote source")
	}

	if extractor.rootNode.Data.Time == "" {
		t.Fatal("attribute `time` did not extracted from `Cube` element")
	}

	for _, cube := range extractor.rootNode.Data.Cubes {
		_, parseErr := strconv.ParseFloat(cube.Rate, 64)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
	}
}

func TestFetchRatesFromLocalFirst(t *testing.T) {
	extractor := NewExtractorXml(sourceLocal)

	if extractor.isExistFile() {
		err := extractor.removeFile()
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	err := extractor.FetchRates()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if extractor.source != sourceRemote {
		t.Fatal("extracted from non-remote source")
	}

	if extractor.rootNode.Data.Time == "" {
		t.Fatal("attribute `time` did not extracted from `Cube` element")
	}

	for _, cube := range extractor.rootNode.Data.Cubes {
		_, parseErr := strconv.ParseFloat(cube.Rate, 64)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
	}
}

func TestFetchRatesFromLocal(t *testing.T) {
	extractor := NewExtractorXml(sourceLocal)

	err := extractor.FetchRates()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !extractor.isExistFile() {
		t.Fatalf("local file `./go-grpc-protobuf/rates.xml` doesn't exist")
	}

	if extractor.source != sourceLocal {
		t.Fatal("extracted from non-local source")
	}

	if extractor.rootNode.Data.Time == "" {
		t.Fatal("attribute `time` did not extracted from `Cube` element")
	}

	for _, cube := range extractor.rootNode.Data.Cubes {
		_, parseErr := strconv.ParseFloat(cube.Rate, 64)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
	}
}

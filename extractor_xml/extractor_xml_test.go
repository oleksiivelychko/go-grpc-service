package extractor_xml

import (
	"strconv"
	"testing"
)

func TestExtractorXML_FetchDataFromURL(t *testing.T) {
	extractorXML := NewExtractorXML(SourceURL)

	err := extractorXML.FetchData()
	if err != nil {
		t.Fatal(err.Error())
	}

	if extractorXML.source != SourceURL {
		t.Fatal("source doesn't equal to SourceURL")
	}

	if extractorXML.RootNode.Data.Time == "" {
		t.Fatal("attribute `time` couldn't extracted from `Cube` element")
	}

	for _, cube := range extractorXML.RootNode.Data.Cubes {
		_, parseFloatErr := strconv.ParseFloat(cube.Rate, 64)
		if parseFloatErr != nil {
			t.Error(parseFloatErr)
		}
	}
}

func TestExtractorXML_FetchDataFromLocalFirstTime(t *testing.T) {
	extractorXML := NewExtractorXML(SourceLocal)

	if extractorXML.isExistFile() {
		err := extractorXML.RemoveFile()
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

	if extractorXML.RootNode.Data.Time == "" {
		t.Fatal("attribute `time` couldn't extracted from `Cube` element")
	}

	for _, cube := range extractorXML.RootNode.Data.Cubes {
		_, parseFloatErr := strconv.ParseFloat(cube.Rate, 64)
		if parseFloatErr != nil {
			t.Error(parseFloatErr)
		}
	}
}

func TestExtractorXML_FetchDataFromLocal(t *testing.T) {
	xmlExtractor := NewExtractorXML(SourceLocal)

	err := xmlExtractor.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !xmlExtractor.isExistFile() {
		t.Fatalf("local file `%s` doesn't exist", localXML)
	}

	if xmlExtractor.source != SourceLocal {
		t.Fatal("source doesn't equal to SourceURL")
	}

	if xmlExtractor.RootNode.Data.Time == "" {
		t.Fatal("attribute `time` couldn't extracted from `Cube` element")
	}

	for _, cube := range xmlExtractor.RootNode.Data.Cubes {
		_, parseFloatErr := strconv.ParseFloat(cube.Rate, 64)
		if parseFloatErr != nil {
			t.Error(parseFloatErr)
		}
	}

	err = xmlExtractor.RemoveFile()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

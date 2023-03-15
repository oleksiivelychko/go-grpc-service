package xml_extractor

import (
	"strconv"
	"testing"
)

func TestFetchDataFromRemote(t *testing.T) {
	xmlExtractor := NewXmlExtractor(SourceRemote)

	err := xmlExtractor.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if xmlExtractor.source != SourceRemote {
		t.Fatal("extracted from non-remote source")
	}

	if xmlExtractor.RootNode.Data.Time == "" {
		t.Fatal("attribute `time` did not extracted from `Cube` element")
	}

	for _, cube := range xmlExtractor.RootNode.Data.Cubes {
		_, parseErr := strconv.ParseFloat(cube.Rate, 64)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
	}
}

func TestFetchDataFromLocalFirstTime(t *testing.T) {
	xmlExtractor := NewXmlExtractor(SourceLocal)

	if xmlExtractor.isExistFile() {
		err := xmlExtractor.removeFile()
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	err := xmlExtractor.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if xmlExtractor.source != SourceRemote {
		t.Fatal("extracted from non-remote source")
	}

	if xmlExtractor.RootNode.Data.Time == "" {
		t.Fatal("attribute `time` did not extracted from `Cube` element")
	}

	for _, cube := range xmlExtractor.RootNode.Data.Cubes {
		_, parseErr := strconv.ParseFloat(cube.Rate, 64)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
	}
}

func TestFetchDataFromLocal(t *testing.T) {
	xmlExtractor := NewXmlExtractor(SourceLocal)

	err := xmlExtractor.FetchData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !xmlExtractor.isExistFile() {
		t.Fatalf("local file `%s` doesn't exist", localXml)
	}

	if xmlExtractor.source != SourceLocal {
		t.Fatal("extracted from non-local source")
	}

	if xmlExtractor.RootNode.Data.Time == "" {
		t.Fatal("attribute `time` did not extracted from `Cube` element")
	}

	for _, cube := range xmlExtractor.RootNode.Data.Cubes {
		_, parseErr := strconv.ParseFloat(cube.Rate, 64)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
	}

	err = xmlExtractor.removeFile()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

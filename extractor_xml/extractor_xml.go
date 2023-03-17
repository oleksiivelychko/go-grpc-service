package extractor_xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const urlXML = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
const localXML = "./rates.xml"

const (
	SourceURL = iota
	SourceLocal
)

type source int8

type ExtractorXML struct {
	RootNode *RootNode
	source   source
}

type RootNode struct {
	Data RootCube `xml:"Cube>Cube"`
}

type RootCube struct {
	Time  string `xml:"time,attr"`
	Cubes []Cube `xml:"Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

func NewExtractorXML(src source) *ExtractorXML {
	return &ExtractorXML{&RootNode{}, src}
}

func (extractorXML *ExtractorXML) FetchData() error {
	if extractorXML.source == SourceURL {
		if err := extractorXML.decodeFromURL(); err != nil {
			return err
		}
	} else {
		if err := extractorXML.readFromLocal(); err != nil {
			return err
		}
	}

	return nil
}

func (extractorXML *ExtractorXML) makeRequest(xmlURL string) (io.ReadCloser, error) {
	resp, err := http.DefaultClient.Get(xmlURL)
	if err != nil {
		return nil, fmt.Errorf("GET `%s` got error: %s", xmlURL, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (extractorXML *ExtractorXML) decodeFromURL() error {
	body, err := extractorXML.makeRequest(urlXML)
	if err != nil {
		return err
	}

	defer body.Close()

	return xml.NewDecoder(body).Decode(extractorXML.RootNode)
}

func (extractorXML *ExtractorXML) readFromURL() ([]byte, error) {
	body, err := extractorXML.makeRequest(urlXML)
	if err != nil {
		return []byte{}, err
	}

	defer body.Close()

	bytes, err := io.ReadAll(body)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to read response body: %v", err)
	}

	extractorXML.source = SourceURL
	return bytes, nil
}

/*
*
readFromLocal tries to read data from local file and if failed - get them from remote source and then save it to file.
*/
func (extractorXML *ExtractorXML) readFromLocal() error {
	var bytes []byte
	var err error

	if !extractorXML.isExistFile() {
		bytes, err = extractorXML.readFromURL()
		if err != nil {
			return err
		}
		_, err = extractorXML.save(bytes)
	} else {
		bytes, err = os.ReadFile(extractorXML.getFilePath())
		if err == nil {
			extractorXML.source = SourceLocal
		}
	}

	return xml.Unmarshal(bytes, &extractorXML.RootNode)
}

/*
*
save writes data (bytes) from remote source into file.
*/
func (extractorXML *ExtractorXML) save(bytes []byte) (int, error) {
	fileXML, err := os.OpenFile(extractorXML.getFilePath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return 0, fmt.Errorf("unable to create `%s` file. %s", localXML, err)
	}
	defer fileXML.Close()

	writtenBytes, err := fileXML.Write(bytes)
	if err != nil {
		return 0, fmt.Errorf("unable to write bytes into file. %s", err)
	}

	return writtenBytes, nil
}

/*
getFilePath returns absolute path to XML file regarding project directory.
*/
func (extractorXML *ExtractorXML) getFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd, localXML)
}

func (extractorXML *ExtractorXML) RemoveFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return os.Remove(filepath.Join(wd, localXML))
}

func (extractorXML *ExtractorXML) isExistFile() bool {
	_, err := os.Stat(extractorXML.getFilePath())
	return !errors.Is(err, os.ErrNotExist)
}

package xml_extractor

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const remoteXml = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
const localXml = "./rates.xml"

const (
	SourceRemote = iota
	SourceLocal
)

type source int8

type XmlExtractor struct {
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

func NewXmlExtractor(src source) *XmlExtractor {
	return &XmlExtractor{&RootNode{}, src}
}

func (xmlExtractor *XmlExtractor) FetchData() error {
	if xmlExtractor.source == SourceRemote {
		if err := xmlExtractor.decodeFromRemote(); err != nil {
			return err
		}
	} else {
		if err := xmlExtractor.readFromLocal(); err != nil {
			return err
		}
	}

	return nil
}

func (xmlExtractor *XmlExtractor) makeRequest(url string) (io.ReadCloser, error) {
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET `%s` got error: %s", url, err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return response.Body, nil
}

func (xmlExtractor *XmlExtractor) decodeFromRemote() error {
	body, err := xmlExtractor.makeRequest(remoteXml)
	if err != nil {
		return err
	}

	defer body.Close()

	return xml.NewDecoder(body).Decode(xmlExtractor.RootNode)
}

func (xmlExtractor *XmlExtractor) readFromRemote() ([]byte, error) {
	body, err := xmlExtractor.makeRequest(remoteXml)
	if err != nil {
		return []byte{}, err
	}

	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return []byte{}, fmt.Errorf("read response body: %v", err)
	}

	xmlExtractor.source = SourceRemote
	return data, nil
}

/*
*
readFromLocal tries to read data from local file and if failed - get them from remote source and then save to file.
*/
func (xmlExtractor *XmlExtractor) readFromLocal() error {
	var data []byte
	var err error

	if !xmlExtractor.isExistFile() {
		data, err = xmlExtractor.readFromRemote()
		if err != nil {
			return err
		}
		_, err = xmlExtractor.save(data)
	} else {
		data, err = os.ReadFile(xmlExtractor.getFilePath())
		if err == nil {
			xmlExtractor.source = SourceLocal
		}
	}

	return xml.Unmarshal(data, &xmlExtractor.RootNode)
}

/*
*
save writes data from remote source into file.
*/
func (xmlExtractor *XmlExtractor) save(bytes []byte) (int, error) {
	xmlFile, err := os.OpenFile(xmlExtractor.getFilePath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return 0, fmt.Errorf("unable to create `%s`. %s", localXml, err)
	}
	defer xmlFile.Close()

	writtenBytes, err := xmlFile.Write(bytes)
	if err != nil {
		return 0, fmt.Errorf("unable to write bytes into file. %s", err)
	}

	return writtenBytes, nil
}

/*
getFilePath returns absolute path to xml file regarding project directory.
*/
func (xmlExtractor *XmlExtractor) getFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd, localXml)
}

func (xmlExtractor *XmlExtractor) removeFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return os.Remove(filepath.Join(wd, localXml))
}

func (xmlExtractor *XmlExtractor) isExistFile() bool {
	_, err := os.Stat(xmlExtractor.getFilePath())
	return !errors.Is(err, os.ErrNotExist)
}

package data

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const xmlUrl = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
const localXml = "./rates.xml"

const (
	SourceRemote = iota
	SourceLocal
)

type source int8

type Extractor struct {
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

func NewExtractor(src source) *Extractor {
	return &Extractor{&RootNode{}, src}
}

func (e *Extractor) FetchRates() error {
	if e.source == SourceRemote {
		if err := e.decodeFromRemote(); err != nil {
			return err
		}
	} else {
		if err := e.readFromLocal(); err != nil {
			return err
		}
	}

	return nil
}

func (e *Extractor) makeRequest(url string) (io.ReadCloser, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET `%s` got error: %s", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (e *Extractor) decodeFromRemote() error {
	body, err := e.makeRequest(xmlUrl)
	if err != nil {
		return err
	}

	defer body.Close()

	return xml.NewDecoder(body).Decode(e.RootNode)
}

func (e *Extractor) readFromRemote() ([]byte, error) {
	body, err := e.makeRequest(xmlUrl)
	if err != nil {
		return []byte{}, err
	}

	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return []byte{}, fmt.Errorf("read response body: %v", err)
	}

	e.source = SourceRemote
	return data, nil
}

/*
*
readFromLocal tries to read data from local file and if failed - get them from remote source and then save to file.
*/
func (e *Extractor) readFromLocal() error {
	var data []byte
	var err error

	if !e.isExistFile() {
		data, err = e.readFromRemote()
		if err != nil {
			return err
		}
		_, err = e.save(data)
	} else {
		data, err = os.ReadFile(e.getFilePath())
		if err == nil {
			e.source = SourceLocal
		}
	}

	return xml.Unmarshal(data, &e.RootNode)
}

/*
*
save writes data from remote source into file.
*/
func (e *Extractor) save(data []byte) (int, error) {
	ratesXml, err := os.OpenFile(e.getFilePath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return 0, fmt.Errorf("unable to create `%s`. %s", localXml, err)
	}
	defer ratesXml.Close()

	writtenBytes, err := ratesXml.Write(data)
	if err != nil {
		return 0, fmt.Errorf("unable to write data into file. %s", err)
	}

	return writtenBytes, nil
}

/*
getFilePath returns absolute path to xml file regarding project directory.
*/
func (e *Extractor) getFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd, localXml)
}

func (e *Extractor) removeFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return os.Remove(filepath.Join(wd, localXml))
}

func (e *Extractor) isExistFile() bool {
	_, err := os.Stat(e.getFilePath())
	return !errors.Is(err, os.ErrNotExist)
}

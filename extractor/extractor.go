package extractor

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const urlXML = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

const (
	SourceURL = iota
	SourceLocal
)

type source int8

type XML struct {
	RootNode *RootNode
	source   source
	localXML string
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

func New(src source, localXML string) *XML {
	return &XML{&RootNode{}, src, localXML}
}

func (extractor *XML) FetchData() error {
	if extractor.source == SourceURL {
		if err := extractor.readFromURL(); err != nil {
			return err
		}
	} else {
		if err := extractor.readFromLocal(); err != nil {
			return err
		}
	}

	return nil
}

func (extractor *XML) readFromURL() error {
	response, err := http.DefaultClient.Get(urlXML)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return xml.NewDecoder(response.Body).Decode(extractor.RootNode)
}

func (extractor *XML) readFromLocal() (err error) {
	var bytesArr []byte

	if _, err = os.Stat(extractor.localXML); !errors.Is(err, os.ErrNotExist) {
		bytesArr, err = os.ReadFile(extractor.localXML)
		if err != nil {
			return err
		}
	} else {
		extractor.source = SourceURL

		bytesArr, err = extractor.readURL()
		if err != nil {
			return err
		}

		err = os.WriteFile(extractor.localXML, bytesArr, 0644)
		if err != nil {
			return err
		}
	}

	return xml.Unmarshal(bytesArr, &extractor.RootNode)
}

func (extractor *XML) readURL() ([]byte, error) {
	resp, err := http.DefaultClient.Get(urlXML)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("got %d status code", resp.StatusCode)
	}

	if err != nil {
		return []byte{}, err
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	bytesArr, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if len(bytesArr) == 0 {
		return []byte{}, fmt.Errorf("response body is empty")
	}

	return bytesArr, nil
}

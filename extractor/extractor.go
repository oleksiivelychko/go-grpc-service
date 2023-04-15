package extractor

import (
	"encoding/xml"
	"github.com/oleksiivelychko/go-grpc-service/utils"
	"github.com/oleksiivelychko/go-grpc-service/utils/request"
	"github.com/oleksiivelychko/go-grpc-service/utils/writer"
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
		if err := extractor.decodeFromURL(); err != nil {
			return err
		}
	} else {
		if err := extractor.readFromLocal(); err != nil {
			return err
		}
	}

	return nil
}

func (extractor *XML) decodeFromURL() error {
	response, err := http.DefaultClient.Get(urlXML)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return xml.NewDecoder(response.Body).Decode(extractor.RootNode)
}

func (extractor *XML) readFromLocal() (err error) {
	var bytesArr []byte

	if utils.IsPathValid(extractor.localXML) {
		bytesArr, err = os.ReadFile(extractor.localXML)
		if err != nil {
			return err
		}
	} else {
		bytesArr, err = request.GET(urlXML)
		if err != nil {
			return err
		}
		extractor.source = SourceURL

		_, err = writer.ToFile(extractor.localXML, bytesArr)
		if err != nil {
			return err
		}
	}

	return xml.Unmarshal(bytesArr, &extractor.RootNode)
}

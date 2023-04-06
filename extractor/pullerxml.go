package extractor

import (
	"encoding/xml"
	"github.com/oleksiivelychko/go-utils/request"
	"github.com/oleksiivelychko/go-utils/system"
	"github.com/oleksiivelychko/go-utils/writer"
	"os"
)

const urlXML = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

const (
	SourceURL = iota
	SourceLocal
)

type source int8

type PullerXML struct {
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

func NewPullerXML(src source, localXML string) *PullerXML {
	return &PullerXML{&RootNode{}, src, localXML}
}

func (puller *PullerXML) FetchData() error {
	if puller.source == SourceURL {
		if err := puller.decodeFromURL(); err != nil {
			return err
		}
	} else {
		if err := puller.readFromLocal(); err != nil {
			return err
		}
	}

	return nil
}

func (puller *PullerXML) decodeFromURL() error {
	body, err := request.DoGET(urlXML)
	if err != nil {
		return err
	}
	defer body.Close()

	return xml.NewDecoder(body).Decode(puller.RootNode)
}

func (puller *PullerXML) readFromLocal() (err error) {
	var bytesArr []byte

	if system.IsPathValid(puller.localXML) {
		bytesArr, err = os.ReadFile(puller.localXML)
		if err != nil {
			return err
		}
	} else {
		bytesArr, err = request.ReadGET(urlXML)
		if err != nil {
			return err
		}
		puller.source = SourceURL

		_, err = writer.ToFile(puller.localXML, bytesArr)
		if err != nil {
			return err
		}
	}

	return xml.Unmarshal(bytesArr, &puller.RootNode)
}

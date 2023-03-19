package extractor_xml

import (
	"encoding/xml"
	"github.com/oleksiivelychko/go-utils/file_ops"
	"github.com/oleksiivelychko/go-utils/request_get"
	"os"
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

func (extractorXML *ExtractorXML) decodeFromURL() error {
	body, err := request_get.DoRequestGET(urlXML)
	if err != nil {
		return err
	}
	defer body.Close()

	return xml.NewDecoder(body).Decode(extractorXML.RootNode)
}

func (extractorXML *ExtractorXML) readFromLocal() (err error) {
	var bytes []byte

	if file_ops.DoesFileExist(localXML) {
		bytes, err = os.ReadFile(localXML)
		if err != nil {
			return err
		}
	} else {
		bytes, err = request_get.DoAndReadRequestGET(urlXML)
		if err != nil {
			return err
		}
		extractorXML.source = SourceURL
		_, err = file_ops.SaveToFile(localXML, bytes)
	}

	return xml.Unmarshal(bytes, &extractorXML.RootNode)
}

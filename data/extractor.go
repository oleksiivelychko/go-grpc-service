package data

import (
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strconv"
)

const targetUrl = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type Exchanger struct {
	log   hclog.Logger
	rates map[string]float64
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

type Cubes struct {
	Data []Cube `xml:"Cube>Cube>Cube"`
}

func NewExchanger(l hclog.Logger) (*Exchanger, error) {
	e := &Exchanger{log: l, rates: map[string]float64{}}
	err := e.fetchRates()
	return e, err
}

func (e *Exchanger) fetchRates() error {
	response, err := http.DefaultClient.Get(targetUrl)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 status code, got %d", response.StatusCode)
	}

	defer response.Body.Close()

	cubes := &Cubes{}
	err = xml.NewDecoder(response.Body).Decode(cubes)
	if err != nil {
		return err
	}

	for _, cube := range cubes.Data {
		rate, parseErr := strconv.ParseFloat(cube.Rate, 64)
		if parseErr != nil {
			e.log.Error("[ERROR] cannot parse the rate value to float", "error", parseErr)
			continue
		}

		e.rates[cube.Currency] = rate
	}

	return nil
}

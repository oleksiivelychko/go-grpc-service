package processor

import (
	"fmt"
	"github.com/oleksiivelychko/go-grpc-service/extractor"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"strconv"
	"time"
)

type Exchanger struct {
	rates     map[string]float64
	pullerXML *extractor.PullerXML
}

func NewExchanger(pullerXML *extractor.PullerXML) (*Exchanger, error) {
	exchanger := &Exchanger{pullerXML: pullerXML, rates: map[string]float64{}}
	err := exchanger.processRates()
	return exchanger, err
}

func (exchanger *Exchanger) GetRate(fromCurrency, toCurrency string) (float64, error) {
	rateFromCurrency, ok := exchanger.rates[fromCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for base [from] %s currency", fromCurrency)
	}

	rateToCurrency, ok := exchanger.rates[toCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for destination [to] %s currency", toCurrency)
	}

	return rateFromCurrency / rateToCurrency, nil
}

func (exchanger *Exchanger) processRates() error {
	err := exchanger.pullerXML.FetchData()
	if err != nil {
		return err
	}

	for _, cube := range exchanger.pullerXML.RootNode.Data.Cubes {
		rate, parseFloatErr := strconv.ParseFloat(cube.Rate, 64)
		if parseFloatErr != nil {
			return fmt.Errorf("unable to parse rate value %s to float: %s", cube.Rate, parseFloatErr)
		}

		exchanger.rates[cube.Currency] = rate
	}

	exchanger.rates[grpcservice.Currencies_EUR.String()] = 1

	return nil
}

/*
TrackRates sends message to the channel when rates are changed.
Fake simulation process, the local development use only.
*/
func (exchanger *Exchanger) TrackRates(interval time.Duration) chan struct{} {
	ch := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				for currency, rate := range exchanger.rates {
					// can be 10% of original value
					change := rand.Float64() / 10
					isPositive := rand.Intn(1)
					if isPositive == 0 {
						// new value will be min 90% of old
						change = 1 - change
					} else {
						// new value will be 110% of old
						change = 1 + change
					}

					exchanger.rates[currency] = rate * change
				}

				// notify updates, this will block unless there is a listener on the other end
				ch <- struct{}{}
			}
		}
	}()

	return ch
}

func (exchanger *Exchanger) GetProtoTimestamp() *timestamppb.Timestamp {
	createdAt, err := time.Parse("2006-01-02", exchanger.pullerXML.RootNode.Data.Time)
	if err != nil {
		createdAt = timestamppb.Now().AsTime()
	}

	return timestamppb.New(createdAt)
}

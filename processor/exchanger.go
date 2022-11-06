package processor

import (
	"fmt"
	"github.com/oleksiivelychko/go-grpc-service/data"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"strconv"
	"time"
)

type Exchanger struct {
	rates     map[string]float64
	extractor *data.Extractor
}

func NewExchanger(e *data.Extractor) (*Exchanger, error) {
	exchanger := &Exchanger{extractor: e, rates: map[string]float64{}}
	err := exchanger.processRates()
	return exchanger, err
}

func (e *Exchanger) GetRate(fromCurrency, toCurrency string) (float64, error) {
	rateFromCurrency, ok := e.rates[fromCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for base [from] currency %s", fromCurrency)
	}

	rateToCurrency, ok := e.rates[toCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for destination [to] currency %s", toCurrency)
	}

	return rateFromCurrency / rateToCurrency, nil
}

func (e *Exchanger) processRates() error {
	err := e.extractor.FetchRates()
	if err != nil {
		return err
	}

	for _, cube := range e.extractor.RootNode.Data.Cubes {
		rate, parseErr := strconv.ParseFloat(cube.Rate, 64)
		if parseErr != nil {
			return fmt.Errorf("cannot parse the rate value `%s` to float. %s", cube.Rate, parseErr)
		}

		e.rates[cube.Currency] = rate
	}

	e.rates[grpc_service.Currencies_EUR.String()] = 1

	return nil
}

/*
TrackRates sends message to the channel when rates are changed.
Fake simulation process, the local development use only.
*/
func (e *Exchanger) TrackRates(interval time.Duration) chan struct{} {
	ch := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				for currency, rate := range e.rates {
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

					e.rates[currency] = rate * change
				}

				// notify updates, this will block unless there is a listener on the other end
				ch <- struct{}{}
			}
		}
	}()

	return ch
}

func (e *Exchanger) GetProtoTime() *timestamppb.Timestamp {
	createdAt, err := time.Parse("2006-01-02", e.extractor.RootNode.Data.Time)
	if err != nil {
		createdAt = timestamppb.Now().AsTime()
	}

	return timestamppb.New(createdAt)
}

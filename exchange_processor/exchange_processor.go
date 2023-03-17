package exchange_processor

import (
	"fmt"
	gs "github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	xe "github.com/oleksiivelychko/go-grpc-service/xml_extractor"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"strconv"
	"time"
)

type ExchangeProcessor struct {
	rates        map[string]float64
	xmlExtractor *xe.XmlExtractor
}

func NewExchangeProcessor(xmlExtractor *xe.XmlExtractor) (*ExchangeProcessor, error) {
	exchangeProcessor := &ExchangeProcessor{xmlExtractor: xmlExtractor, rates: map[string]float64{}}
	err := exchangeProcessor.processRates()
	return exchangeProcessor, err
}

func (exchangeProcessor *ExchangeProcessor) GetRate(fromCurrency, toCurrency string) (float64, error) {
	rateFromCurrency, ok := exchangeProcessor.rates[fromCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for base [from] currency %s", fromCurrency)
	}

	rateToCurrency, ok := exchangeProcessor.rates[toCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for destination [to] currency %s", toCurrency)
	}

	return rateFromCurrency / rateToCurrency, nil
}

func (exchangeProcessor *ExchangeProcessor) processRates() error {
	err := exchangeProcessor.xmlExtractor.FetchData()
	if err != nil {
		return err
	}

	for _, cube := range exchangeProcessor.xmlExtractor.RootNode.Data.Cubes {
		rate, parseErr := strconv.ParseFloat(cube.Rate, 64)
		if parseErr != nil {
			return fmt.Errorf("cannot parse the rate value `%s` to float. %s", cube.Rate, parseErr)
		}

		exchangeProcessor.rates[cube.Currency] = rate
	}

	exchangeProcessor.rates[gs.Currencies_EUR.String()] = 1

	return nil
}

/*
TrackRates sends message to the channel when rates are changed.
Fake simulation process, the local development use only.
*/
func (exchangeProcessor *ExchangeProcessor) TrackRates(interval time.Duration) chan struct{} {
	ch := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				for currency, rate := range exchangeProcessor.rates {
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

					exchangeProcessor.rates[currency] = rate * change
				}

				// notify updates, this will block unless there is a listener on the other end
				ch <- struct{}{}
			}
		}
	}()

	return ch
}

func (exchangeProcessor *ExchangeProcessor) GetProtoTimestamp() *timestamppb.Timestamp {
	createdAt, err := time.Parse("2006-01-02", exchangeProcessor.xmlExtractor.RootNode.Data.Time)
	if err != nil {
		createdAt = timestamppb.Now().AsTime()
	}

	return timestamppb.New(createdAt)
}

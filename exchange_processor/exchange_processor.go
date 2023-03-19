package exchange_processor

import (
	"fmt"
	"github.com/oleksiivelychko/go-grpc-service/extractor_xml"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"strconv"
	"time"
)

type ExchangeProcessor struct {
	rates        map[string]float64
	extractorXML *extractor_xml.ExtractorXML
}

func NewExchangeProcessor(extractorXML *extractor_xml.ExtractorXML) (*ExchangeProcessor, error) {
	exchangeProcessor := &ExchangeProcessor{extractorXML: extractorXML, rates: map[string]float64{}}
	err := exchangeProcessor.processRates()
	return exchangeProcessor, err
}

func (exchangeProcessor *ExchangeProcessor) GetRate(fromCurrency, toCurrency string) (float64, error) {
	rateFromCurrency, ok := exchangeProcessor.rates[fromCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for base [from] '%s' currency", fromCurrency)
	}

	rateToCurrency, ok := exchangeProcessor.rates[toCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for destination [to] '%s' currency", toCurrency)
	}

	return rateFromCurrency / rateToCurrency, nil
}

func (exchangeProcessor *ExchangeProcessor) processRates() error {
	err := exchangeProcessor.extractorXML.FetchData()
	if err != nil {
		return err
	}

	for _, cube := range exchangeProcessor.extractorXML.RootNode.Data.Cubes {
		rate, parseFloatErr := strconv.ParseFloat(cube.Rate, 64)
		if parseFloatErr != nil {
			return fmt.Errorf("unable to parse rate value '%s' to float: %s", cube.Rate, parseFloatErr)
		}

		exchangeProcessor.rates[cube.Currency] = rate
	}

	exchangeProcessor.rates[grpc_service.Currencies_EUR.String()] = 1

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
	createdAt, err := time.Parse("2006-01-02", exchangeProcessor.extractorXML.RootNode.Data.Time)
	if err != nil {
		createdAt = timestamppb.Now().AsTime()
	}

	return timestamppb.New(createdAt)
}

package exchanger

import (
	"fmt"
	"github.com/oleksiivelychko/go-grpc-service/extractor"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"strconv"
	"time"
)

type Processor struct {
	rates  map[string]float64
	puller *extractor.PullerXML
}

func NewProcessor(pullerXML *extractor.PullerXML) (*Processor, error) {
	processor := &Processor{puller: pullerXML, rates: map[string]float64{}}
	err := processor.processRates()
	return processor, err
}

func (processor *Processor) GetRate(fromCurrency, toCurrency string) (float64, error) {
	rateFromCurrency, ok := processor.rates[fromCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for base [from] %s currency", fromCurrency)
	}

	rateToCurrency, ok := processor.rates[toCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for destination [to] %s currency", toCurrency)
	}

	return rateFromCurrency / rateToCurrency, nil
}

func (processor *Processor) processRates() error {
	err := processor.puller.FetchData()
	if err != nil {
		return err
	}

	for _, cube := range processor.puller.RootNode.Data.Cubes {
		rate, parseFloatErr := strconv.ParseFloat(cube.Rate, 64)
		if parseFloatErr != nil {
			return fmt.Errorf("unable to parse rate value %s to float: %s", cube.Rate, parseFloatErr)
		}

		processor.rates[cube.Currency] = rate
	}

	processor.rates[grpcservice.Currencies_EUR.String()] = 1

	return nil
}

/*
TrackRates sends message to the channel when rates are changed.
Fake simulation process, the local development use only.
*/
func (processor *Processor) TrackRates(interval time.Duration) chan struct{} {
	ch := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				for currency, rate := range processor.rates {
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

					processor.rates[currency] = rate * change
				}

				// notify updates, this will block unless there is a listener on the other end
				ch <- struct{}{}
			}
		}
	}()

	return ch
}

func (processor *Processor) GetProtoTimestamp() *timestamppb.Timestamp {
	createdAt, err := time.Parse("2006-01-02", processor.puller.RootNode.Data.Time)
	if err != nil {
		createdAt = timestamppb.Now().AsTime()
	}

	return timestamppb.New(createdAt)
}

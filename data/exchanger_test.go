package data

import (
	"fmt"
	"testing"
)

func TestNewExchanger(t *testing.T) {
	e, err := NewExchanger()

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("Rates: %#v", e.rates)
}

package data

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"testing"
)

func TestNewExchanger(t *testing.T) {
	e, err := NewExchanger(hclog.Default())

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("Rates: %#v", e.rates)
}

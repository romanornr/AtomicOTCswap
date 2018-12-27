package bcoins

import (
	"testing"
)

var viacoin = Coin{}

func TestBtcutil_GetBtcUtil(t *testing.T) {
	viacoin.Symbol = "via"

	input, _ := NewAmount(0.001)
	expected := "0.001 VIA"

	if input.String() != expected {
		t.Errorf("error expected %v but got %v instead\n", expected, input)
	}
}

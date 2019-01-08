package bcoins

import (
	"fmt"
	btcutil "github.com/viacoin/viautil"
	"testing"
)

var viacoin = Coin{}

func TestBtcutil_GetBtcUtil(t *testing.T) {
	viacoin.Symbol = "xzc"

	input, _ := btcutil.NewAmount(0.001)
	expected := "0.001 VIA"

	fmt.Println(input) // should be 0.001 VIA not 1000

	if input.String() != expected {
		t.Errorf("error expected %v but got %v instead\n", expected, input)
	}
}

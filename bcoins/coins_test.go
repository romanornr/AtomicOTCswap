package bcoins

import (
	"github.com/viacoin/viautil"
	"testing"
)

var viacoin = Coin{}
func TestBtcutil_GetBtcUtil(t *testing.T) {
	viacoin.Symbol = "via"

	input := viacoin.GetBtcUtil((2200000))

	expected := viautil.Amount(2200000)

	if input != expected.String() {
		t.Errorf("error expected %v but got %v instead\n", expected, input)
	}
}

package bcoins

import (
	"fmt"
	"github.com/viacoin/viautil"
	"testing"
)

var viacoin = Coin{}
func TestBtcutil_GetBtcUtil(t *testing.T) {
	viacoin.Symbol = "via"

	input := AmountUnit(22)

	expected := viautil.Amount(2200000)
	fmt.Println(expected)

	if input.String() != expected.String() {
		t.Errorf("error expected %v but got %v instead\n", expected, input)
	}
}

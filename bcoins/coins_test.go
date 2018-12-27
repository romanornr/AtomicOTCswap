package bcoins

import (
	//"fmt"
	//btcutil "github.com/viacoin/viautil"
	"fmt"
	"testing"
)

var viacoin = Coin{}
func TestBtcutil_GetBtcUtil(t *testing.T) {
	viacoin.Symbol = "via"


	input, _ := NewAmount(0.001)

	fmt.Println(input)

	//expected, _ := btcutil.NewAmount(0.001)

	//if input.String() != expected.String() {
	//	t.Errorf("error expected %v but got %v instead\n", expected, input)
	//}
}

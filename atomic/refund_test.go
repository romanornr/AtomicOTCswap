package atomic

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	btcutil "github.com/viacoin/viautil"
	"testing"
)

func TestRefund(t *testing.T) {
	wif, _ := btcutil.DecodeWIF("WVtSmyCEuUNeXgPEnC4uEPJPRswpCjS2uUCe6wykRaYEBre9Y4fX")
	a, err := Refund("6382012088a8202820d90f5829d64f84e0d1b3bf97e75a464e49ff13da503512bd694add023eb58876a91424cc424c1e5e977175d2b20012554d39024bd68f6704d680215cb17576a9148e4d174eb2236d8ab0c5bb250df08a5af054d3206888ac", "02000000015fcb00efe16361bb0685ae10937ab13fc48e89bb8a4a3b55c6b294822620d77c000000006b4830450221008fdd8d99ec6824a9f229ab5a8f5b8b4365c6bf46f369e81a4627c1b4bd37d172022011298f8f1f63a46f423d41da89fee48ef5daa37e8079c5f4bd92b6330d6780f0012102a7b08bb2a3609223a185761231d815e287ec13b74ccff3feb274253f7737356affffffff01a08601000000000017a914d8010216698392f7224dc43639f5dd69698a58238700000000", wif)
	if err != nil {
		fmt.Println(err)
	}

	coin, _ := bcoins.SelectCoin("via")

	err = a.Run(wif, &coin)
	if err != nil {
		fmt.Println(err)
	}

}

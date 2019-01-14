package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	btcutil "github.com/viacoin/viautil"
)

func main() {
	contractHex := "6382012088a820a6f25c9323b65cd9bd61991d332f484eeb0a57f5c5355ce5b5326b2abc0e09908876a914885635ae37c7de6d85cf06684a66e9921dcdf872670448643d5cb17576a9148e4d174eb2236d8ab0c5bb250df08a5af054d3206888ac"
	transactionHex := "0200000001cc9d027bad81a789e02f0a99dee668696f51d82ca8c81545b641b98fec4a458b000000006b483045022100bbd2f7138e5d22dae7bdc27e614a2ef2607c13b1955131f5d37c7f7ae7e3646302204dc9568b2dfdba5edb423714181385f6065b2402c609484900fb2634026554c7012102a7b08bb2a3609223a185761231d815e287ec13b74ccff3feb274253f7737356affffffff01a0c44a000000000017a9140ef4216f7cfb9c0f64b5d5e3905fbc226a670d548700000000"
	wif, _ := btcutil.DecodeWIF("WVtSmyCEuUNeXgPEnC4uEPJPRswpCjS2uUCe6wykRaYEBre9Y4fX")
	coin, _ := bcoins.SelectCoin("via")

	err := atomic.Refund(contractHex, transactionHex, wif, coin)
	if err != nil {
		fmt.Println(err)
	}
}

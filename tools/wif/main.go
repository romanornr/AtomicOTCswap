package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viautil"
)

func main() {
	//wif, _ := atomic.GenerateNewWIF()
	coin, _ := bcoins.SelectCoin("ltc")
	chaincfg.Register(coin.Network.ChainCgfMainNetParams())

	wif, _ := viautil.DecodeWIF("6voC8Xn5jfd4oNw3sDCFhwR3SearCzjVHcQ8sHHnHPPggEawMSS")

	pk, _ := atomic.GenerateNewPublicKey(*wif, &coin)
	fmt.Printf("public key: %s\n", pk.AddressPubKeyHash().String())
}

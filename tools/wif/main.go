package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
)

func main() {
	wif, _ := atomic.GenerateNewWIF()

	pk, _ := atomic.GenerateNewPublicKey(*wif)
	fmt.Printf("public key: %s\n WIF: %s\n", pk.AddressPubKeyHash(), wif.PrivKey)
}

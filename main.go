package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"github.com/viacoin/viautil"
)
func main() {
	//initiatorWIF, _ := atomic.GenerateNewWIF()
	partyBAddress := "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A"

	//fmt.Println(initiatorWIF.String())

	wif, _ := viautil.DecodeWIF("WVtSmyCEuUNeXgPEnC4uEPJPRswpCjS2uUCe6wykRaYEBre9Y4fX")
	pk, _ := atomic.GenerateNewPublicKey(*wif)
	fmt.Printf("public key to deposit 0.001 on: %s\n", pk.AddressPubKeyHash())

	err := atomic.Initiate(partyBAddress, wif, 0.001)

	if err != nil {
		fmt.Printf("%s", err)
	}
	//fmt.Println(a.amount)
}

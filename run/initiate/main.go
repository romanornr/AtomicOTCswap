package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"github.com/viacoin/viautil"
)

func main() {
	//initiatorWIF, _ := atomic.GenerateNewWIF()
	//partyBAddress := "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A"

	partyBAddress := "LS8MLtQfz4nmDHeaEFd2rMBZCyywmRRchj"

	//fmt.Println(initiatorWIF.String())

	wif, _ := viautil.DecodeWIF("6vC5Y2CuYgW4vMJ1D2iSKWP8zTyAcPHopxAwpF2BYPqR7ioFAzQ")

	//coin, _ := bcoins.SelectCoin("ltc")
	//pk, _ := atomic.GenerateNewPublicKey(*wif, &coin)
	//fmt.Printf("public key to deposit 0.001 on: %s\n", pk.AddressPubKeyHash())

	err := atomic.Initiate("ltc", partyBAddress, wif, 0.05)

	if err != nil {
		fmt.Printf("%s", err)
	}
	//fmt.Println(a.amount)
}

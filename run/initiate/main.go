package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	btcutil "github.com/viacoin/viautil"
)

func main() {
	//initiatorWIF, _ := atomic.GenerateNewWIF()
	partyBAddress := "VnRi5kfigWPgDAz62mY4oLKRruErZwTZJB"

	//partyBAddress := "LS8MLtQfz4nmDHeaEFd2rMBZCyywmRRchj"

	//fmt.Println(initiatorWIF.String())
	wif, _ := btcutil.DecodeWIF("WVtSmyCEuUNeXgPEnC4uEPJPRswpCjS2uUCe6wykRaYEBre9Y4fX")

	coin, _ := bcoins.SelectCoin("via")
	pk, _ := atomic.GenerateNewPublicKey(*wif, &coin)
	fmt.Printf("public key to deposit 0.001 on: %s\n", pk.AddressPubKeyHash())

	err := atomic.Initiate("via", partyBAddress, wif, 0.049) //TODO fix fee with inputs

	if err != nil {
		fmt.Printf("%s", err)
	}
	//fmt.Println(a.amount)
}

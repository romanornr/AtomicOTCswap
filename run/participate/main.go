package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"github.com/viacoin/viautil"
)

func main() {

	//reader := bufio.NewReader(os.Stdin)

	//initiatorWIF, _ := atomic.GenerateNewWIF()
	partyBAddress := "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A"

	//fmt.Println(initiatorWIF.String())

	wif, _ := viautil.DecodeWIF("WVtSmyCEuUNeXgPEnC4uEPJPRswpCjS2uUCe6wykRaYEBre9Y4fX")
	pk, _ := atomic.GenerateNewPublicKey(*wif)
	fmt.Printf("public key to deposit 0.001 on: %s\n", pk.AddressPubKeyHash())

	fmt.Print("Enter secret hash: ")
	//secret, _:= reader.ReadString('\n')
	secret := "636a9caa905eb7d673a9d55e3af2b7e60c12d1f866b7773e2c5c50eca910d6d8"
	err := atomic.Participate("via", partyBAddress, wif, 0.05, secret)

	if err != nil {
		fmt.Printf("%s", err)
	}
	//fmt.Println(a.amount)
}

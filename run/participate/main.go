package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	btcutil "github.com/viacoin/viautil"
	"log"
)

func main() {

	//reader := bufio.NewReader(os.Stdin)

	//initiatorWIF, _ := atomic.GenerateNewWIF()
	partyBAddress := "LiLCgk1wHc2zBNp1ZxQs1PX7NUGy1Jx4tu"

	//fmt.Println(initiatorWIF.String())

	//wif, _ := viautil.DecodeWIF("WVtSmyCEuUNeXgPEnC4uEPJPRswpCjS2uUCe6wykRaYEBre9Y4fX")
	//pk, _ := atomic.GenerateNewPublicKey(*wif)
	//fmt.Printf("public key to deposit 0.001 on: %s\n", pk.AddressPubKeyHash())
	wif, err := btcutil.DecodeWIF("6vCjF7aNQyvJHyNS3JhYTvNR4TEZDTxwUQRhyKewm5pcTNzzbGQ")
	if err != nil {
		log.Panicf("%s\n", err)
	}
	coin, _ := bcoins.SelectCoin("ltc")

	pk, _ := atomic.GenerateNewPublicKey(*wif, &coin)
	fmt.Println(pk.AddressPubKeyHash())
	//fmt.Print("Enter secret hash: ")
	//secret, _:= reader.ReadString('\n')
	secret := "a6f25c9323b65cd9bd61991d332f484eeb0a57f5c5355ce5b5326b2abc0e0990"
	err = atomic.Participate("ltc", partyBAddress, wif, 0.05, secret)

	if err != nil {
		fmt.Printf("%s", err)
	}
	//fmt.Println(a.amount)
}

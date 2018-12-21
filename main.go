package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
)
func main() {
	//initiatorRefundWIF, _ := atomic.GenerateNewWIF()
	//initiatorRefundAddr, _ := atomic.GenerateNewPublicKey(*initiatorRefundWIF)
	partyBAddress := "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A"


	err := atomic.Initiate(partyBAddress, "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A", 2.09)
	if err != nil {
		fmt.Printf("%s", err)
	}
	//fmt.Println(a.amount)
}

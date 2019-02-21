package key

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"strings"
	"testing"
)

func TestGetAtomicSwapAddresses(t *testing.T) {
	depositAsset, _ := bcoins.SelectCoin("via")
	receivingAsset, _ := bcoins.SelectCoin("ltc")

	addresses, _ := GetAtomicSwapAddresses(&depositAsset, &receivingAsset)

	if strings.HasPrefix(addresses.DepositAddress, "V") == false {
		fmt.Println("fail")
	}
	fmt.Println(addresses)
}

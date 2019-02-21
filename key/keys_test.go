package key

import (
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"strings"
	"testing"
)

var tests = []struct {
	assetSymbol   string
	addressPrefix string
	wifPrefix     string
}{
	{assetSymbol: "via", addressPrefix: "V", wifPrefix: "W"}, // viacoin
	{assetSymbol: "ltc", addressPrefix: "L", wifPrefix: "T"}, // litecoin
}

// test if address and privatekey/wif have the right prefixes
// these tests are for compressed addresses and compressed WIF keys.
func TestSwapAddresses(t *testing.T) {

	for _, pair := range tests {
		asset, _ := bcoins.SelectCoin(pair.assetSymbol)
		net := asset.Network.ChainCgfMainNetParams()
		wif, _ := GenerateNewWIF(net)
		address, _ := GenerateNewPublicKey(*wif, net)

		// test WIF keys
		if strings.HasPrefix(wif.String(), pair.wifPrefix) != true {
			t.Error(
				"For", asset.Name,
				"expected", pair.wifPrefix,
				"got", wif.String(),
			)
		}

		// test addresses
		if strings.HasPrefix(address.EncodeAddress(), pair.addressPrefix) != true {
			t.Error(
				"For", asset.Name,
				"expected", pair.wifPrefix,
				"got", wif.String(),
			)
		}
	}
}

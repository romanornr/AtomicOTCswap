package key

import (
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/btcec"
	"github.com/viacoin/viad/chaincfg"
	btcutil "github.com/viacoin/viautil"
)

func GenerateNewWIF(coin *bcoins.Coin) (*btcutil.WIF, error) {
	chaincfg.Register(coin.Network.ChainCgfMainNetParams())
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}

	wif, err := btcutil.NewWIF(secret, coin.Network.ChainCgfMainNetParams(), true)
	return wif, err
}

func GenerateNewPublicKey(wif btcutil.WIF, coin *bcoins.Coin) (*btcutil.AddressPubKey, error) {
	chaincfg.Register(coin.Network.ChainCgfMainNetParams())
	pk, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), coin.Network.ChainCgfMainNetParams())
	return pk, err
}

type SwapAddresses struct {
	DepositAddress   string `json:"deposit_address"`
	ReceivingAddress string `json:"receiving_address"`
	DepositWif       string `json:"deposit_wif"`
	ReceivingWif     string `json:"receiving_wif"`
}

func GetAtomicSwapAddresses(depositAsset *bcoins.Coin, receivingAsset *bcoins.Coin) (swapAddresses SwapAddresses, err error) {
	depositWif, err := GenerateNewWIF(depositAsset)
	if err != nil {
		return swapAddresses, err
	}

	receivingWif, err := GenerateNewWIF(receivingAsset)
	if err != nil {
		return swapAddresses, err
	}

	depositAddress, err := GenerateNewPublicKey(*depositWif, depositAsset)
	if err != nil {
		return swapAddresses, err
	}

	receivingAddress, err := GenerateNewPublicKey(*receivingWif, receivingAsset)
	if err != nil {
		return swapAddresses, err
	}

	return SwapAddresses{
		DepositAddress:   depositAddress.EncodeAddress(),
		DepositWif:       depositWif.String(),
		ReceivingAddress: receivingAddress.EncodeAddress(),
		ReceivingWif:     receivingWif.String(),
	}, nil
}

package swaputil

import (
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/btcec"
	"github.com/viacoin/viad/chaincfg"
	btcutil "github.com/viacoin/viautil"
)

func GenerateNewWIF(net *chaincfg.Params) (*btcutil.WIF, error) {
	chaincfg.Register(net)
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}

	wif, err := btcutil.NewWIF(secret, net, true)
	return wif, err
}

func GenerateNewPublicKey(wif btcutil.WIF, net *chaincfg.Params) (*btcutil.AddressPubKey, error) {
	chaincfg.Register(net)
	pk, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), net)
	return pk, err
}

type SwapAddresses struct {
	DepositAddress   string `json:"deposit_address"`
	ReceivingAddress string `json:"receiving_address"`
	DepositWif       string `json:"deposit_wif"`
	ReceivingWif     string `json:"receiving_wif"`
}

func GetAtomicSwapAddresses(depositAsset *bcoins.Coin, receivingAsset *bcoins.Coin) (swapAddresses SwapAddresses, err error) {
	depositWif, err := GenerateNewWIF(depositAsset.Network.ChainCgfMainNetParams())
	if err != nil {
		return swapAddresses, err
	}

	receivingWif, err := GenerateNewWIF(receivingAsset.Network.ChainCgfMainNetParams())
	if err != nil {
		return swapAddresses, err
	}

	depositAddress, err := GenerateNewPublicKey(*depositWif, depositAsset.Network.ChainCgfMainNetParams())
	if err != nil {
		return swapAddresses, err
	}

	receivingAddress, err := GenerateNewPublicKey(*receivingWif, receivingAsset.Network.ChainCgfMainNetParams())
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

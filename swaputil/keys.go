// Copyright (c) 2019 Romano (Viacoin developer)
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

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

type SwapKeyPair struct {
	DepositAddress   string `json:"deposit_address"`
	DepositWif       string `json:"deposit_wif"`
	ReceivingAddress string `json:"receiving_address"`
	ReceivingWif     string `json:"receiving_wif"`
}

// generate key pair for the deposit assset and receiving asset. Return a SwapKeyPair struct
func GenerateSwapKeyPair(depositAsset *bcoins.Coin, receivingAsset *bcoins.Coin) (swapKeyPair SwapKeyPair, err error) {
	depositWif, err := GenerateNewWIF(depositAsset.Network.ChainCgfMainNetParams())
	if err != nil {
		return swapKeyPair, err
	}

	receivingWif, err := GenerateNewWIF(receivingAsset.Network.ChainCgfMainNetParams())
	if err != nil {
		return swapKeyPair, err
	}

	depositAddress, err := GenerateNewPublicKey(*depositWif, depositAsset.Network.ChainCgfMainNetParams())
	if err != nil {
		return swapKeyPair, err
	}

	receivingAddress, err := GenerateNewPublicKey(*receivingWif, receivingAsset.Network.ChainCgfMainNetParams())
	if err != nil {
		return swapKeyPair, err
	}

	return SwapKeyPair{
		DepositAddress:   depositAddress.EncodeAddress(),
		DepositWif:       depositWif.String(),
		ReceivingAddress: receivingAddress.EncodeAddress(),
		ReceivingWif:     receivingWif.String(),
	}, nil
}

// Copyright (c) 2019 Romano (Viacoin developer)
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package bcoins

import (
	"fmt"
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viad/wire"
	"strings"
)

type Coin struct {
	Name          string
	Unit          string
	Symbol        string
	Network       *Network
	Insight       *Insight
	TxVersion     int32
	MinRelayTxFee float64
	FeePerByte    int
	Dust          int64
}

type Insight struct {
	Explorer string
	Api      string
}

type Network struct {
	Name           string
	P2PKH          byte
	P2SH           byte
	PrivateKeyID   byte
	HDCoinType     uint32
	HDPublicKeyID  [4]byte
	HDPrivateKeyID [4]byte
	magic          wire.BitcoinNet
}

var coins = map[string]Coin{
	"via": {Name: "viacoin", Symbol: "via", Unit: "VIA", Network: &Network{Name: "viacoin", P2PKH: 0x47, P2SH: 0x21, PrivateKeyID: 0xC7, HDCoinType: 14, magic: 0xcbc6680f},
		Insight: &Insight{"https://explorer.viacoin.org", "https://explorer.viacoin.org/api"}, TxVersion: 2, MinRelayTxFee: 0.001, FeePerByte: 110, Dust: int64(1000),
	},
	"btc": {Name: "bitcoin", Symbol: "btc", Network: &Network{Name:"bitcoin", P2PKH: 0x00, P2SH:0x05, PrivateKeyID: 0x80, HDCoinType: 1, HDPublicKeyID: [4]byte{0x04, 0x88, 0xb2, 0x1e}, HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, magic: 0xfbc0b6db},
		Insight: &Insight{"https://insight.bitpay.com", "https://insight.bitpay.com/api"}, TxVersion: 2, MinRelayTxFee: 0.00001, FeePerByte: 13, Dust: int64(546),
	},

	"btct": {Name: "bitcoin testnet", Symbol: "btct", Network: &Network{Name:"bitcoin testnet", P2PKH: 0x6f, P2SH:0xc4, PrivateKeyID: 0xef, HDCoinType: 0, HDPublicKeyID: [4]byte{0x04, 0x35, 0x87, 0xcf}, HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, magic: 0x0709110b}, //testnet 3
		Insight: &Insight{"https://insight.bitpay.com", "https://test-insight.bitpay.com/api"}, TxVersion: 2, MinRelayTxFee: 0.00001, FeePerByte: 13, Dust: int64(546),
	},

	"ltc": {Name: "litecoin", Symbol: "ltc", Unit: "LTC",
		Network: &Network{Name: "litecoin", P2PKH: 0x30, P2SH: 0x32, PrivateKeyID: 0xB0, HDCoinType: 2, HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, HDPublicKeyID: [4]byte{0x04, 0x88, 0xb2, 0x1e}, magic: 0xfbc0b6db},
		Insight: &Insight{"https://insight.litecore.io", "https://insight.litecore.io/api"}, TxVersion: 2, MinRelayTxFee: 0.001, FeePerByte: 280, Dust: int64(10000),
	},

	//TODO add Decred, Zcoin, Vercoin
}

// select a coin by symbol and return Coin struct and error
// coin symbol to lower case
func SelectCoin(symbol string) (Coin, error) {
	if coins, ok := coins[strings.ToLower(symbol)]; ok {
		return coins, nil
	}
	return Coin{}, fmt.Errorf("altcoin %s not found\n", symbol)
}

// set the chainparams correct for the given Network struct
// and returns the chaincfg.Params
func (network Network) ChainCgfMainNetParams() *chaincfg.Params {
	networkParams := &chaincfg.MainNetParams
	networkParams.Name = network.Name
	networkParams.Net = network.magic
	networkParams.PubKeyHashAddrID = network.P2PKH
	networkParams.ScriptHashAddrID = network.P2SH
	networkParams.HDCoinType = network.HDCoinType
	networkParams.HDPrivateKeyID = network.HDPrivateKeyID
	networkParams.HDPublicKeyID = network.HDPublicKeyID
	networkParams.PrivateKeyID = network.PrivateKeyID
	return networkParams
}

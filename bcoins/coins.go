package bcoins

import (
	"fmt"
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viad/wire"
)

type Coin struct {
	Symbol string
	Name string
	Network *Network
	Insight *Insight
	TxVersion int32
}

type Insight struct {
	Explorer string
	Api string
}

type Network struct {
	Name string
	P2PKH byte
	P2SH byte
	PrivateKeyID byte
	magic wire.BitcoinNet
}

var coins = map[string]Coin {
	"via": {Name: "viacoin", Symbol: "via", Network: &Network{"viacoin", 0x47,0x21,  0xC7, 0xcbc6680f},
		Insight: &Insight{"https://explorer.viacoin.org", "https://explorer.viacoin.org/api"}, TxVersion: 2,
	},
}

// select a coin by symbol and return Coin struct and error
func SelectCoin(symbol string) (Coin, error) {
	if coins, ok := coins[symbol]; ok {
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
	networkParams.PrivateKeyID = network.PrivateKeyID
	return networkParams
}
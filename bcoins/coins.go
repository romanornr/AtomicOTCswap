package bcoins

import (
	"fmt"
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viad/wire"
	"strconv"
)

type Coin struct {
	Symbol string
	Name string
	Network *Network
	Insight *Insight
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
		Insight: &Insight{"https://explorer.viacoin.org", "https://explorer.viacoin.org/api"},
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
func (network Network) GetNetworkParams() *chaincfg.Params {
	networkParams := &chaincfg.MainNetParams
	networkParams.Name = network.Name
	networkParams.Net = network.magic
	networkParams.PubKeyHashAddrID = network.P2PKH
	networkParams.ScriptHashAddrID = network.P2SH
	networkParams.PrivateKeyID = network.PrivateKeyID
	return networkParams
}



type AmountUnit int

// These constants define various units used when describing a viacoin
// monetary amount.
const (
	AmountMegaBTC  AmountUnit = 6
	AmountKiloBTC  AmountUnit = 3
	AmountBTC      AmountUnit = 0
	AmountMilliBTC AmountUnit = -3
	AmountMicroBTC AmountUnit = -6
	AmountSatoshi  AmountUnit = -8
)

func (coin Coin) GetBtcUtil(u AmountUnit) string {
	//u := coin.Btcutil.AmountUnit
	//u := btcutil.AmountUnit

	// String returns the unit as a string.  For recognized units, the SI
	// prefix is used, or "Satoshi" for the base unit.  For all unrecognized
	// units, "1eN BTC" is returned, where N is the AmountUnit.
	//func (u AmountUnit) String() string {
		switch u {
	case AmountMegaBTC:
		return "MVIA"
	case AmountKiloBTC:
		return "kVIA"
	case AmountBTC:
		return "VIA"
	case AmountMilliBTC:
		return "mVIA"
	case AmountMicroBTC:
		return "μVIA"
	case AmountSatoshi:
		return "Satoshi"
	default:
		return "1e" + strconv.FormatInt(int64(u), 10) + " VIA"
	}

}

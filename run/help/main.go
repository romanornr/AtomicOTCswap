package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viad/wire"
	btcutil "github.com/viacoin/viautil"
)

type Coin struct {
	Symbol        string
	Name          string
	Network       *Network
	Insight       *Insight
	TxVersion     int32
	MinRelayTxFee float64
	FeePerByte    int
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
	HDPrivateKeyID [4]byte
	HDPublicKeyID  [4]byte
	magic          wire.BitcoinNet
}

var coins = map[string]Coin{
	"via": {Name: "viacoin", Symbol: "via", Network: &Network{Name: "viacoin", P2PKH: 0x47, P2SH: 0x21, PrivateKeyID: 0xC7, HDCoinType: 14, magic: 0xcbc6680f},
		Insight: &Insight{"https://explorer.viacoin.org", "https://explorer.viacoin.org/api"}, TxVersion: 2, MinRelayTxFee: 0.001, FeePerByte: 110,
	},

	"ltc": {Name: "litecoin", Symbol: "ltc",
		Network: &Network{Name: "litecoin", P2PKH: 0x30, P2SH: 0x32, PrivateKeyID: 0xB0, HDCoinType: 2, HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, HDPublicKeyID: [4]byte{0x04, 0x88, 0xb2, 0x1e}, magic: 0xfbc0b6db},
		Insight: &Insight{"https://insight.litecore.io", "https://insight.litecore.io/api"}, TxVersion: 2, MinRelayTxFee: 0.001, FeePerByte: 280,
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
	networkParams.HDCoinType = network.HDCoinType
	networkParams.HDPrivateKeyID = network.HDPrivateKeyID
	networkParams.HDPublicKeyID = network.HDPublicKeyID
	networkParams.PrivateKeyID = network.PrivateKeyID
	return networkParams
}

func init() {

}

func main() {

	coin, err := bcoins.SelectCoin("ltc")
	if err != nil {
		fmt.Println("err")
	}

	//partyBAddress := "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A"  // viacoin address does work
	partyBAddress := "LS8MLtQfz4nmDHeaEFd2rMBZCyywmRRchj"
	fmt.Println(coin.Network.ChainCgfMainNetParams().PubKeyHashAddrID)

	counterParty2Addr, err := btcutil.DecodeAddress(partyBAddress, coin.Network.ChainCgfMainNetParams())
	if err != nil {
		fmt.Printf("failed to decode the address from the participant: %s\n", err)
	}

	counterParty2AddrP2KH, ok := counterParty2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		fmt.Println("participant address is not P2KH")
	}

	fmt.Println(counterParty2Addr)
	fmt.Println(counterParty2AddrP2KH)

}


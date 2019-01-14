package atomic

import (
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/btcec"
	"github.com/viacoin/viad/chaincfg"
	btcutil "github.com/viacoin/viautil"
)

func GenerateNewWIF() (*btcutil.WIF, error){
	secret , err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}

	wif, err := btcutil.NewWIF(secret, &chaincfg.MainNetParams, true)
	return wif, err
}

func GenerateNewPublicKey(wif btcutil.WIF, coin *bcoins.Coin) (*btcutil.AddressPubKey, error) {
	chaincfg.Register(coin.Network.ChainCgfMainNetParams())
	pk, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), coin.Network.ChainCgfMainNetParams())
	return pk, err
}

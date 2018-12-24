package atomic

import (
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

func GenerateNewPublicKey(wif btcutil.WIF) (*btcutil.AddressPubKey, error){
	pk, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.MainNetParams)
	return pk, err
}

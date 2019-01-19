package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/chaincfg"
)

func main() {
	//coin, _ := bcoins.SelectCoin("ltc")
	coin, _ := bcoins.SelectCoin("ltc")
	chaincfg.Register(coin.Network.ChainCgfMainNetParams())
	contractHex := "6382012088a8209a06fa2f0ce4ee01156fcaf260d535a71a279fd144de073611b939d6d93076548876a914fd796daea971996803ef49a13312ed7e52dcc8db6704182f435cb17576a914b0a89c190a64ad1d23aec98f34459cb1577437876888ac"
	contractTransactionHex := "020000000170ef4077ef5ad34af31913b5f1ac9a68e83a0a00ef42ae8a2bfe992335f98652000000006b483045022100b0d76d0f9264855fba4d404c82e173de25802c85a0a7cf38eef93293d6b4c6f8022069d1665ecd6f69d1dde17927a1c6f54d3acb5851887978cc8646db789a12eee00121037f98c2ee049eff3868f9361968c9e2abf4cbfc6be81c536be86ac034202a2c83ffffffff01003e49000000000017a914fd4ceb40417276a477e979c02e6fd18f4a365a438700000000"
	secret := "9a06fa2f0ce4ee01156fcaf260d535a71a279fd144de073611b939d6d9307654"
	wif := "6vCjF7aNQyvJHyNS3JhYTvNR4TEZDTxwUQRhyKewm5pcTNzzbGQ"
	//wif, err := btcutil.DecodeWIF("WVtSmyCEuUNeXgPEnC4uEPJPRswpCjS2uUCe6wykRaYEBre9Y4fX")


	redemption, err := atomic.Redeem(coin.Symbol, contractHex, contractTransactionHex, secret, wif)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(redemption.TransactionHex)
}

package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/chaincfg"
	btcutil"github.com/viacoin/viautil"
)

func main() {
	//coin, _ := bcoins.SelectCoin("ltc")
	contractHex := "6382012088a820a6f25c9323b65cd9bd61991d332f484eeb0a57f5c5355ce5b5326b2abc0e09908876a914fd796daea971996803ef49a13312ed7e52dcc8db6704f5643d5cb17576a914b0a89c190a64ad1d23aec98f34459cb1577437876888ac"
	//contractHex := "231f7fde2859217b4f7649a127d77a417ebc5f59e80c98914c921c70d160caac"
	//contractTransactionHex := "0200000001590672b850060af486de4b097afb91a7e4c17261dbfccc0bdaa348a6dd7188c8000000006b48304502210098e629a643134a49f4a8c917379997578a8a80dc9067aba829ddda31d45074a002201fd2a2bbdbcf157daee4b0dab9b31e0d475b8b0dabce7a9698f125c35fffda790121037f98c2ee049eff3868f9361968c9e2abf4cbfc6be81c536be86ac034202a2c83ffffffff01404b4c000000000017a91481e525114aa96f548ac60f4f8183732c3789b5ac8700000000"
	contractTransactionHex := "0200000001590672b850060af486de4b097afb91a7e4c17261dbfccc0bdaa348a6dd7188c8000000006b48304502210098e629a643134a49f4a8c917379997578a8a80dc9067aba829ddda31d45074a002201fd2a2bbdbcf157daee4b0dab9b31e0d475b8b0dabce7a9698f125c35fffda790121037f98c2ee049eff3868f9361968c9e2abf4cbfc6be81c536be86ac034202a2c83ffffffff01404b4c000000000017a91481e525114aa96f548ac60f4f8183732c3789b5ac8700000000"
	secret := "75a9c26148aa6eb40a7040c0ffb1cb2327d682ef635b3f14e29444c06369dc99"
	//wif, err := btcutil.DecodeWIF("WVtSmyCEuUNeXgPEnC4uEPJPRswpCjS2uUCe6wykRaYEBre9Y4fX")
	wif, err := btcutil.DecodeWIF("6vCjF7aNQyvJHyNS3JhYTvNR4TEZDTxwUQRhyKewm5pcTNzzbGQ")
	coin, _ := bcoins.SelectCoin("ltc")
	chaincfg.Register(coin.Network.ChainCgfMainNetParams())

	err = atomic.Redeem(coin.Symbol, contractHex, contractTransactionHex, secret, wif)
	if err != nil {
		fmt.Println(err)
	}
}

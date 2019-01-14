package main

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/atomic"
)

func main() {
	//contract, _ := atomic.AuditContract("ltc", "6382012088a820a6f25c9323b65cd9bd61991d332f484eeb0a57f5c5355ce5b5326b2abc0e09908876a914fd796daea971996803ef49a13312ed7e52dcc8db6704f5643d5cb17576a914b0a89c190a64ad1d23aec98f34459cb1577437876888ac", "0200000001590672b850060af486de4b097afb91a7e4c17261dbfccc0bdaa348a6dd7188c8000000006b48304502210098e629a643134a49f4a8c917379997578a8a80dc9067aba829ddda31d45074a002201fd2a2bbdbcf157daee4b0dab9b31e0d475b8b0dabce7a9698f125c35fffda790121037f98c2ee049eff3868f9361968c9e2abf4cbfc6be81c536be86ac034202a2c83ffffffff01404b4c000000000017a91481e525114aa96f548ac60f4f8183732c3789b5ac8700000000")
	//fmt.Println(contract.Show())

	contract, _ := atomic.AuditContract("via", "6382012088a820a6f25c9323b65cd9bd61991d332f484eeb0a57f5c5355ce5b5326b2abc0e09908876a914fd796daea971996803ef49a13312ed7e52dcc8db6704f5643d5cb17576a914b0a89c190a64ad1d23aec98f34459cb1577437876888ac", "0200000001590672b850060af486de4b097afb91a7e4c17261dbfccc0bdaa348a6dd7188c8000000006b48304502210098e629a643134a49f4a8c917379997578a8a80dc9067aba829ddda31d45074a002201fd2a2bbdbcf157daee4b0dab9b31e0d475b8b0dabce7a9698f125c35fffda790121037f98c2ee049eff3868f9361968c9e2abf4cbfc6be81c536be86ac034202a2c83ffffffff01404b4c000000000017a91481e525114aa96f548ac60f4f8183732c3789b5ac8700000000")
	fmt.Println(contract.Show())
}




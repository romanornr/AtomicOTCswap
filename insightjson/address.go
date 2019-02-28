// Copyright (c) 2017 The Decred developers
// Copyright (c) 2019 Romano (Viacoin developer)
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package insightjson

type UnspentOutputs []struct {
	Address       string  `json:"address"`
	Txid          string  `json:"txid"`
	Vout          int     `json:"vout"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Satoshis      int     `json:"satoshis"`
	Height        int     `json:"height"`
	Confirmations int     `json:"confirmations"`
}

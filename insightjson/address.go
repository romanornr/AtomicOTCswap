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

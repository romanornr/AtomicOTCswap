package insightjson


package insightjson

type Tx struct {
	Txid          string  `json:"txid,omitempty"`
	Version       int32   `json:"version,omitempty"`
	Locktime      uint32  `json:"locktime"`
	Vins          []*Vin  `json:"vin,omitempty"`
	Vouts         []*Vout `json:"vout,omitempty"`
	Blockhash     string  `json:"blockhash,omitempty"`
	Blockheight   int64   `json:"blockheight"`
	Confirmations uint64  `json:"confirmations"`
	Time          int64   `json:"time,omitempty"`
	Blocktime     int64   `json:"blocktime,omitempty"`
	IsCoinBase    bool    `json:"isCoinBase,omitempty"`
	ValueOut      float64 `json:"valueOut"`
	Size          uint32  `json:"size,omitempty"`
	ValueIn       float64 `json:"valueIn,omitempty"`
	Fees          float64 `json:"fees,omitempty"`
}

type Vin struct {
	Txid      string     `json:"txid,omitempty"`
	Vout      uint32     `json:"vout"`
	Sequence  uint32     `json:"sequence,omitempty"`
	N         int        `json:"n"`
	ScriptSig *ScriptSig `json:"scriptSig,omitempty"`
	Addr      string     `json:"addr,omitempty"`
	ValueSat  int64      `json:"valueSat"`
	Value     float64    `json:"value,omitempty"`
	CoinBase  string     `json:"coinbase,omitempty"`
}

type Vout struct {
	Value        string       `json:"value"`
	N            uint32       `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey,omitempty"`
	SpentTxID    interface{}  `bson:"spentTxId" json:"spentTxId"`
	SpentIndex   interface{}  `bson:"spentIndex" json:"spentIndex"`
	SpentHeight  interface{}  `bson:"spentHeight" json:"spentHeight"`
}

type ScriptPubKey struct {
	Hex       string   `json:"hex,omitempty"`
	Asm       string   `json:"asm,omitempty"`
	Addresses []string `json:"addresses,omitempty"`
	Type      string   `json:"type,omitempty"`
}

type ScriptSig struct {
	Hex string `json:"hex,omitempty"`
	Asm string `json:"asm,omitempty"`
}

type InsightRawTx struct {
	Rawtx string `json:"rawtx"`
}

type Txid struct {
	Txid string `json:"txid"`
}

package bcoins

type Transaction struct {
	AssetName          string `json:"asset_name"`
	AssetSymbol        string `json:"asset_symbol"`
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

package insight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/romanornr/AtomicOTCswap/insightjson"
	"io/ioutil"
	"net/http"
)

//type errorBroadcast string
const ErrNotEnoughBalance = "16: bad-txns-vout-negative. Code:-26"
const ErrNotEnoughFee = "66: insufficient priority. Code:-26"
const ErrTransactionTooSmall = "64: dust. Code:-26"
const ErrTxDecodeFailed = "Something seems wrong: TX decode failed. Code:-2"

//broadcast a signed transaction to the blockexplorer.
//the transaction will either be denied or rejected by the network
func BroadcastTransaction(asset bcoins.Coin, tx bcoins.Transaction) (insightjson.Txid, bcoins.Transaction, error) {
	jsonData := insightjson.InsightRawTx{Rawtx: tx.SignedTx}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("error broadcasting %s tx because jsonMarshal failed: %s\n", asset.Name, err)
	}

	insightExplorer, _ := GetInsightExplorer(asset.Symbol)
	insightExplorerBroacastApi := fmt.Sprintf("%s/tx/send", insightExplorer.Api)

	response, err := http.Post(insightExplorerBroacastApi, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("Error broadcasting %s transaction with blockexplorer: %s\n", asset.Name, err)
	}

	defer response.Body.Close()

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("error reading response from %s blockexplorer broadcast: %s\n", asset.Name, err)
	}

	if response.StatusCode != 200 { // some error handling if broadcasting fails
		rejectReason := string(result)
		switch rejectReason {
		case ErrNotEnoughBalance:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("not enough balance to cover the transaction including the required fees")
		case ErrNotEnoughFee:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("fee needs to be higher")
		case ErrTransactionTooSmall:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("transaction too small (dust transaction)\nTx does not meet the minimal amount")
		case ErrTxDecodeFailed:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("transaction decode failed !\n Maybe a wrong address?")
		default:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("%s\n", string(result))
		}
	}

	var txid = insightjson.Txid{}
	err = json.Unmarshal([]byte(result), &txid)
	if err != nil {
		return txid, bcoins.Transaction{}, fmt.Errorf("something went wrong with receiving your txid")
	}

	return txid, tx, nil
}

// Copyright (c) 2019 Romano (Viacoin developer)
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package insight

import (
	"encoding/json"
	"fmt"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/romanornr/AtomicOTCswap/insightjson"
	"github.com/viacoin/viad/chaincfg/chainhash"
	"log"
	"net/http"
	"sort"
	"time"
)

type UTXO struct {
	Hash      *chainhash.Hash
	TxIndex   uint32
	Amount    int64
	Spendable bool
	PKScript  []byte
}

func GetUnspentOutputs(publicKey string, coin *bcoins.Coin) []*UTXO {
	insightExplorer := coin.Insight.Api
	url := fmt.Sprintf("%s/addr/%s/utxo", insightExplorer, publicKey)
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("NewRequest error: %s\n", err)
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error client.DO: %s\n", err)
	}
	defer resp.Body.Close()

	var utxos insightjson.UnspentOutputs
	if err := json.NewDecoder(resp.Body).Decode(&utxos); err != nil {
		log.Fatalf("error decoding viacoin insightjson utxo: %s\n", err)
	}

	var sourceUTXOs []*UTXO

	// todo optimize UTXO
	for _, data := range utxos {
		var utxo UTXO
		utxo.Amount = int64(data.Satoshis)
		utxo.TxIndex = uint32(data.Vout)
		hash, _ := chainhash.NewHashFromStr(data.Txid)
		utxo.Hash = hash
		utxo.Spendable = true
		sourceUTXOs = append(sourceUTXOs, &utxo)
	}

	return sourceUTXOs
}

//take in requiredAmount in satoshi. This is how much we want to spend
//however we don't want to use all utxo's we have on our address
//only the minimal utxo's combined is what we want so the transaction becomes smaller in bytes
//the smaller the transaction, the less fees is required
func getMinimalRequiredUTXOByFirstInFirstOut(requiredAmount int64, sourceUTXOs []*UTXO) []*UTXO {
	var newUTXOSet []*UTXO
	var amountSatInUTXOset int64

	for idx := range sourceUTXOs {
		amountSatInUTXOset += sourceUTXOs[idx].Amount
		newUTXOSet = append(newUTXOSet, sourceUTXOs[idx])

		if amountSatInUTXOset > requiredAmount { // break if we have enough inputs to pay the amount + fees
			break
		}
	}
	return newUTXOSet
}

//take in requiredAmount in satoshi. This is how much we want to spend
//however we don't want to use all utxo's we have on our address
//only the minimal utxo's combined is what we want so the transaction becomes smaller in bytes
//the smaller the transaction, the less fees is required
func GetMinimalRequiredUTXO(requiredAmount int64, sourceUTXOs []*UTXO) []*UTXO {
	var newUTXOSet []*UTXO
	var amountSatInUTXOset int64

	//sort the UTXO from low to high
	sort.SliceStable(sourceUTXOs, func(i, j int) bool {
		return sourceUTXOs[i].Amount < sourceUTXOs[j].Amount
	})

	//check if 1 UTXO is enough
	for idx := range sourceUTXOs {
		if sourceUTXOs[idx].Amount >= requiredAmount {
			newUTXOSet = append(newUTXOSet, sourceUTXOs[idx])
			return newUTXOSet
		}
	}

	//high to low and use the highest amounts first until it is equal or more than the required amount
	//here we basically combine inputs
	for idx := range sourceUTXOs {
		idx = len(sourceUTXOs) - 1 - idx // reverse range from high amounts to low amounts
		amountSatInUTXOset += sourceUTXOs[idx].Amount
		newUTXOSet = append(newUTXOSet, sourceUTXOs[idx])

		if amountSatInUTXOset > requiredAmount { // break if we have enough inputs to pay the amount + fees
			break
		}
	}

	return newUTXOSet
}

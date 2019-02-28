// Copyright (c) 2019 Romano (Viacoin developer)
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package insight

import (
	"fmt"
	"github.com/romanornr/AtomicOTCswap/bcoins"
)

func GetInsightExplorer(symbol string) (bcoins.Insight, error) {
	coin, err := bcoins.SelectCoin(symbol)
	if err != nil {
		return bcoins.Insight{}, fmt.Errorf("this altcoin does not have an insight explorer")
	}
	return *coin.Insight, nil
}

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
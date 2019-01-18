package atomic

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/chaincfg"
	btcutil "github.com/viacoin/viautil"
	"strings"
	"time"
)

type participateCmd struct {
	counterParty1Addr *btcutil.AddressPubKeyHash
	amount            btcutil.Amount
	secretHash        []byte
}

type ParticipatedContract struct {
	Coin                   string
	Unit                   string
	ContractAmount         float64
	ContractFee            float64
	ContractRefundFee      float64
	ContractAddress        string
	ContractHex            string
	ContractTransactionID  string
	ContractTransactionHex string
	RefundTransactionID    string
	RefundTransaction      string
}

func Participate(coinTicker string, participantAddr string, wif *btcutil.WIF, amount float64, secret string) (contract ParticipatedContract, err error) {

	coin, err := bcoins.SelectCoin(coinTicker)
	if err != nil {
		return contract, err
	}

	chaincfg.Register(coin.Network.ChainCgfMainNetParams())

	counterParty1Addr, err := btcutil.DecodeAddress(participantAddr, coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return contract, fmt.Errorf("failed to decode the address from the participant: %s", err)
	}

	counterParty1AddrP2KH, ok := counterParty1Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return contract, errors.New("participant address is not P2KH")
	}

	amount2, err := btcutil.NewAmount(amount)
	if err != nil {
		return contract, err
	}

	secretHash, err := hex.DecodeString(secret)
	if err != nil {
		return contract, errors.New("secret hash must be hex encoded")
	}

	cmd := &participateCmd{counterParty1Addr: counterParty1AddrP2KH, amount: amount2, secretHash: secretHash}
	return cmd.runCommand(wif, &coin, amount)
}

func (cmd *participateCmd) runCommand(wif *btcutil.WIF, coin *bcoins.Coin, amount float64) (contract ParticipatedContract, err error) {

	locktime := time.Now().Add(12 * time.Hour).Unix()

	build, err := buildContract(&contractArgs{
		coin1:      coin,
		them:       cmd.counterParty1Addr,
		amount:     cmd.amount,
		locktime:   locktime,
		secretHash: cmd.secretHash,
	}, wif)
	if err != nil {
		return contract, err
	}

	unit := strings.ToUpper(coin.Symbol)
	refundTxHash := build.refundTx.TxHash()

	var contractBuf bytes.Buffer
	contractBuf.Grow(build.contractTx.SerializeSize())
	build.contractTx.Serialize(&contractBuf)

	var refundBuf bytes.Buffer
	refundBuf.Grow(build.refundTx.SerializeSize())
	build.refundTx.Serialize(&refundBuf)

	contract = ParticipatedContract{
		Coin: coin.Name,
		Unit: unit,

		ContractAmount:    amount,
		ContractFee:       build.contractFee.ToBTC(),
		ContractRefundFee: build.refundFee.ToBTC(),
		ContractAddress:   fmt.Sprintf("%v", build.contractP2SH),

		ContractTransactionID:  fmt.Sprintf("%x", build.contractTxHash),
		ContractTransactionHex: fmt.Sprintf("%x", contractBuf.Bytes()),

		RefundTransactionID: fmt.Sprintf("%v", &refundTxHash),
		RefundTransaction:   fmt.Sprintf("%x", refundBuf.Bytes()),
	}

	return contract, err
}

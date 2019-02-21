package atomic

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/chaincfg"
	btcutil "github.com/viacoin/viautil"
	"log"
	"time"
)

type participateCmd struct {
	counterParty1Addr *btcutil.AddressPubKeyHash
	amount            btcutil.Amount
	secretHash        []byte
}

type ParticipatedContract struct {
	Coin                   string  `json:"coin"`
	Unit                   string  `json:"unit"`
	ContractAmount         float64 `json:"contract_amount"`
	ContractFee            float64 `json:"contract_fee"`
	ContractRefundFee      float64 `json:"contract_refund_fee"`
	CounterPartyAddress    string  `json:"counter_party_address"`
	ContractAddress        string  `json:"contract_address"`
	ContractHex            string  `json:"contract_hex"`
	ContractTransactionID  string  `json:"contract_transaction_id"`
	ContractTransactionHex string  `json:"contract_transaction_hex"`
	RefundTransactionID    string  `json:"refund_transaction_id"`
	RefundTransaction      string  `json:"refund_transaction"`
}

func Participate(coinTicker string, participantAddr string, WIFstring string, amount float64, secretHash string) (contract ParticipatedContract, err error) {

	coin, err := bcoins.SelectCoin(coinTicker)
	if err != nil {
		return contract, err
	}

	chaincfg.Register(coin.Network.ChainCgfMainNetParams())

	wif, err := btcutil.DecodeWIF(WIFstring)
	if err != nil {
		log.Printf("error decoding private key in wif format: %s\n", err)
	}

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

	secret, err := hex.DecodeString(secretHash)
	if err != nil {
		return contract, errors.New("secret hash must be hex encoded")
	}

	cmd := &participateCmd{counterParty1Addr: counterParty1AddrP2KH, amount: amount2, secretHash: secret}
	return cmd.runCommand(wif, &coin, amount)
}

func (cmd *participateCmd) runCommand(wif *btcutil.WIF, coin *bcoins.Coin, amount float64) (contract ParticipatedContract, err error) {

	locktime := time.Now().Add(1 * time.Hour).Unix()

	build, err := buildContract(&contractArgs{
		coin:       coin,
		them:       cmd.counterParty1Addr,
		amount:     cmd.amount,
		locktime:   locktime,
		secretHash: cmd.secretHash,
	}, wif)
	if err != nil {
		return contract, err
	}

	refundTxHash := build.refundTx.TxHash()

	var contractBuf bytes.Buffer
	contractBuf.Grow(build.contractTx.SerializeSize())
	build.contractTx.Serialize(&contractBuf)

	var refundBuf bytes.Buffer
	refundBuf.Grow(build.refundTx.SerializeSize())
	build.refundTx.Serialize(&refundBuf)

	contract = ParticipatedContract{
		Coin: coin.Name,
		Unit: coin.Unit,

		ContractAmount:    amount,
		ContractFee:       build.contractFee.ToBTC(),
		ContractRefundFee: build.refundFee.ToBTC(),

		CounterPartyAddress: fmt.Sprintf("%s", cmd.counterParty1Addr.String()),
		ContractAddress:     fmt.Sprintf("%v", build.contractP2SH),
		ContractHex:         fmt.Sprintf("%x", build.contract),

		ContractTransactionID:  fmt.Sprintf("%x", build.contractTxHash),
		ContractTransactionHex: fmt.Sprintf("%x", contractBuf.Bytes()),

		RefundTransactionID: fmt.Sprintf("%v", &refundTxHash),
		RefundTransaction:   fmt.Sprintf("%x", refundBuf.Bytes()),
	}

	return contract, err
}

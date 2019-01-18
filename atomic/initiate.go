package atomic

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/chaincfg"
	btcutil "github.com/viacoin/viautil"
	"strings"
	"time"
)

type initiateCmd struct {
	counterParty2Addr *btcutil.AddressPubKeyHash
	amount            btcutil.Amount
}

type InitiatedContract struct {
	Coin string
	Unit string
	ContractAmount float64
	ContractFee float64
	ContractRefundFee float64
	ContractAddress string
	ContractHex string
	ContractTransactionID string
	ContractTransactionHex string
	RefundTransactionID string
	RefundTransaction string
	Secret string
	SecretHash string
}

func Initiate(coinTicker string, participantAddr string, amount float64, wif *btcutil.WIF) (InitiatedContract, error) {

	coin, err := bcoins.SelectCoin(coinTicker)
	if err != nil {
		return InitiatedContract{}, err
	}

	chaincfg.Register(coin.Network.ChainCgfMainNetParams())

	counterParty2Addr, err := btcutil.DecodeAddress(participantAddr, coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return InitiatedContract{}, fmt.Errorf("failed to decode the address from the participant: %s\n", err)
	}

	counterParty2AddrP2KH, ok := counterParty2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return InitiatedContract{}, errors.New("participant address is not P2KH")
	}

	amount2, err := btcutil.NewAmount(amount)
	if err != nil {
		return InitiatedContract{}, err
	}

	cmd := &initiateCmd{counterParty2Addr: counterParty2AddrP2KH, amount: amount2}
	return cmd.runCommand(wif, &coin, amount)
}

func (cmd *initiateCmd) runCommand(wif *btcutil.WIF, coin *bcoins.Coin, amount float64) (InitiatedContract, error) {
	var secret [secretSize]byte
	_, err := rand.Read(secret[:])
	if err != nil {
		return InitiatedContract{}, err
	}

	secretHash := sha256Hash(secret[:])
	locktime := time.Now().Add(12 * time.Hour).Unix() // NEED TO CHANGE

	build, err := buildContract(&contractArgs{
		coin1:      coin,
		them:       cmd.counterParty2Addr,
		amount:     cmd.amount,
		locktime:   locktime,
		secretHash: secretHash,
	}, wif)

	if err != nil {
		return InitiatedContract{}, err
	}

	unit := strings.ToUpper(coin.Symbol)
	refundTxHash := build.refundTx.TxHash()

	var contractBuf bytes.Buffer
	contractBuf.Grow(build.contractTx.SerializeSize())
	build.contractTx.Serialize(&contractBuf)

	var refundBuf bytes.Buffer
	refundBuf.Grow(build.refundTx.SerializeSize())
	build.refundTx.Serialize(&refundBuf)

	contract := InitiatedContract{

		Secret: fmt.Sprintf("%x", secret),
		SecretHash: fmt.Sprintf("%x", secretHash),

		Coin: coin.Name,
		Unit: unit,

		ContractAmount:amount,
		ContractFee: build.contractFee.ToBTC(),
		ContractRefundFee: build.refundFee.ToBTC(),

		ContractAddress: fmt.Sprintf("%v", build.contractP2SH),
		ContractHex: fmt.Sprintf("%x", build.contract),

		ContractTransactionID: fmt.Sprintf("%x", build.contractTxHash),
		ContractTransactionHex: fmt.Sprintf("%x", contractBuf.Bytes()),

		RefundTransactionID: fmt.Sprintf("%v", &refundTxHash),
		RefundTransaction: fmt.Sprintf("%x", refundBuf.Bytes()),
	}
	return contract, nil
}

package atomic

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	btcutil "github.com/viacoin/viautil"
	"strings"
	"time"
)

type participateCmd struct {
	counterParty1Addr *btcutil.AddressPubKeyHash
	amount            btcutil.Amount
	secretHash        []byte
}

func Participate(coinTicker string, participantAddr string, wif *btcutil.WIF, amount float64, secret string) error {

	coin, err := bcoins.SelectCoin(coinTicker)
	if err != nil {
		return err
	}

	counterParty1Addr, err := btcutil.DecodeAddress(participantAddr, coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return fmt.Errorf("failed to decode the address from the participant: %s", err)
	}

	counterParty1AddrP2KH, ok := counterParty1Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return errors.New("participant address is not P2KH")
	}

	amount2, err := btcutil.NewAmount(amount)
	if err != nil {
		return err
	}

	secretHash, err := hex.DecodeString(secret)
	if err != nil {
		return errors.New("secret hash must be hex encoded")
	}

	cmd := &participateCmd{counterParty1Addr: counterParty1AddrP2KH, amount: amount2, secretHash:secretHash}
	return cmd.runCommand(wif, &coin, amount)
}

func (cmd *participateCmd) runCommand(wif *btcutil.WIF, coin *bcoins.Coin, amount float64) error {


	locktime := time.Now().Add(10 * time.Minute).Unix()

	build, err := buildContract(&contractArgs{
		coin1:      coin,
		them:       cmd.counterParty1Addr,
		amount:     cmd.amount,
		locktime:   locktime,
		secretHash: cmd.secretHash,
	}, wif)
	if err != nil {
		return err
	}

	ticker := strings.ToUpper(coin.Symbol)
	refundTxHash := build.refundTx.TxHash()

	fmt.Printf("Contract fee: %v %s\n", build.contractFee, ticker)
	fmt.Printf("Refund fee:   %v %s\n\n", build.refundFee, ticker)
	fmt.Printf("Contract (%v):\n", build.contractP2SH)
	fmt.Printf("%x\n\n", build.contract)
	var contractBuf bytes.Buffer
	contractBuf.Grow(build.contractTx.SerializeSize())
	build.contractTx.Serialize(&contractBuf)
	fmt.Printf("Contract transaction (%v):\n", build.contractTxHash)
	fmt.Printf("%x\n\n", contractBuf.Bytes())
	var refundBuf bytes.Buffer
	refundBuf.Grow(build.refundTx.SerializeSize())
	build.refundTx.Serialize(&refundBuf)
	fmt.Printf("Refund transaction (%v):\n", &refundTxHash)
	fmt.Printf("%x\n\n", refundBuf.Bytes())

	//build.refundTx.TxHash()
	//calcFeePerKb(build.contractFee, build.contractTx.SerializeSize())
	//calcFeePerKb(build.refundFee, build.refundTx.SerializeSize())
	//
	//var contractBuf bytes.Buffer
	//contractBuf.Grow(build.contractTx.SerializeSize())
	//build.contractTx.Serialize(&contractBuf)
	//
	//var refundBuf bytes.Buffer
	//refundBuf.Grow(build.refundTx.SerializeSize())
	//build.refundTx.Serialize(&refundBuf)

	return nil
}

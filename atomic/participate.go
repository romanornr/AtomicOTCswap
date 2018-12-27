package atomic

import (
	"bytes"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/viacoin/viad/chaincfg"
	btcutil "github.com/viacoin/viautil"
	"time"
)

type participateCmd struct {
	counterparty1Addr *btcutil.AddressPubKeyHash
	amount            btcutil.Amount
	secretHash        []byte
}

func Participate(initiatorAddr string, wif *btcutil.WIF, amount float64) error {
	counterParty1Addr, err := btcutil.DecodeAddress(initiatorAddr, &chaincfg.MainNetParams)
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

	cmd := &participateCmd{counterparty1Addr: counterParty1AddrP2KH, amount: amount2}
	return cmd.runCommand(wif)
}

func (cmd *participateCmd) runCommand(wif *btcutil.WIF) error {

	locktime := time.Now().Add(10 * time.Minute).Unix()

	build, err := buildContract(&contractArgs{
		them:       cmd.counterparty1Addr,
		amount:     cmd.amount,
		locktime:   locktime,
		secretHash: cmd.secretHash,
	}, wif)
	if err != nil {
		return err
	}

	//refundTxHash := build.refundTx.TxHash()
	//contractFeePerKb := calcFeePerKb(build.contractFee, build.contractTx.SerializeSize())
	//refundFeePerKb := calcFeePerKb(build.refundFee, build.refundTx.SerializeSize())

	build.refundTx.TxHash()
	calcFeePerKb(build.contractFee, build.contractTx.SerializeSize())
	calcFeePerKb(build.refundFee, build.refundTx.SerializeSize())

	var contractBuf bytes.Buffer
	contractBuf.Grow(build.contractTx.SerializeSize())
	build.contractTx.Serialize(&contractBuf)

	var refundBuf bytes.Buffer
	refundBuf.Grow(build.refundTx.SerializeSize())
	build.refundTx.Serialize(&refundBuf)

	return nil
}

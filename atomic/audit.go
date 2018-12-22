package atomic

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-errors/errors"
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viad/txscript"
	"github.com/viacoin/viad/wire"
	btcutil "github.com/viacoin/viautil"
)

type auditContractCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

func AuditContract(contractHex string, contractTransaction string) {
	contract, err := hex.DecodeString(contractHex)
	if err != nil {
		fmt.Errorf("failed to decode contract: %v\n", err)
	}

	contractTxBytes, err := hex.DecodeString(contractTransaction)
	if err != nil {
		fmt.Errorf("failed to decode transaction:%v\n", err)
	}
	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		fmt.Errorf("failed to decode transaction: %v\n", err)
	}

	cmd := &auditContractCmd{contract: contract, contractTx: &contractTx}
	cmd.runCommand()
}

func (cmd *auditContractCmd) runCommand() error {
	contractHash160 := btcutil.Hash160(cmd.contract)
	contractOut := -1
	for i, out := range cmd.contractTx.TxOut {
		sc, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.MainNetParams)
		if err != nil || sc != txscript.ScriptHashTy {
			continue
		}
		if bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash160) {
			contractOut = i
		}
	}

	if contractOut == -1 {
		return errors.New("transaction does not contain the contract output")
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, cmd.contract)
	if err != nil {
		return err
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool")
	}
	if pushes.SecretSize != secretSize {
		return fmt.Errorf("contract specifies strange range secret size: %v\n", pushes.SecretSize)
	}

	contractAddr, err := btcutil.NewAddressScriptHash(cmd.contract, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}

	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:], &chaincfg.MainNetParams)
	if err != nil {
		return err
	}
	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], &chaincfg.MainNetParams)
	if err != nil {
		return err
	}

	fmt.Printf("Contract address:        %v\n", contractAddr)
	fmt.Printf("Contract value:          %v\n", btcutil.Amount(cmd.contractTx.TxOut[contractOut].Value))
	fmt.Printf("Recipient address:       %v\n", recipientAddr)
	fmt.Printf("Recipient refund address: %v\n\n", refundAddr)

	if pushes.LockTime >= int64(txscript.LockTimeThreshold) {
		t := time.Unix(pushes.LockTime, 0)
		fmt.Printf("Locktime: %v\n", t.UTC())
		reachedAt := time.Until(t).Truncate(time.Second)
		if reachedAt > 0 {
			fmt.Printf("Locktime reached in %v\n", reachedAt)
		} else {
			fmt.Printf("Contract refund time lock has expired !\n")
			return nil
		}
		fmt.Printf("Locktime: block %v\n", pushes.LockTime)
	}
	return err
}

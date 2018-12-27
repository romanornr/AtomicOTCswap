package atomic

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/viacoin/viad/txscript"
	"github.com/viacoin/viad/wire"
	btcutil "github.com/viacoin/viautil"
	"time"
)

type AuditContractCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

type AuditedContract struct {
	Address                *btcutil.AddressScriptHash `json:"address"`
	Value                  btcutil.Amount             `json:"value"`
	RecipientAddress       *btcutil.AddressPubKeyHash `json:"recipient_address"`
	RecipientRefundAddress *btcutil.AddressPubKeyHash `json:"recipient_refund_address"`
	LockTime               int64                      `json:"lock_time"`
	AtomicSwapDataPushes   *txscript.AtomicSwapDataPushes
}

func AuditContract(coinTicker string, contractHex string, contractTransaction string) (AuditedContract, error) {
	coin, err := bcoins.SelectCoin(coinTicker)
	if err != nil {
		return AuditedContract{}, err
	}

	contract, err := hex.DecodeString(contractHex)
	if err != nil {
		return AuditedContract{}, fmt.Errorf("failed to decode contract: %v\n", err)
	}

	contractTxBytes, err := hex.DecodeString(contractTransaction)
	if err != nil {
		return AuditedContract{}, fmt.Errorf("failed to decode transaction:%v\n", err)
	}
	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return AuditedContract{}, fmt.Errorf("failed to decode transaction: %v\n", err)
	}

	c := AuditContractCmd{contract: contract, contractTx: &contractTx}
	return c.runAudit(coin)
}

func (cmd *AuditContractCmd) runAudit(coin bcoins.Coin) (AuditedContract, error) {
	contractHash160 := btcutil.Hash160(cmd.contract)
	contractOut := -1
	for i, out := range cmd.contractTx.TxOut {
		sc, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, coin.Network.ChainCgfMainNetParams())
		if err != nil || sc != txscript.ScriptHashTy {
			continue
		}
		if bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash160) {
			contractOut = i
		}
	}

	if contractOut == -1 {
		return AuditedContract{}, errors.New("transaction does not contain the contract output")
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, cmd.contract)
	if err != nil {
		return AuditedContract{}, err
	}
	if pushes == nil {
		return AuditedContract{}, errors.New("contract is not an atomic swap script recognized by this tool")
	}
	if pushes.SecretSize != secretSize {
		return AuditedContract{}, fmt.Errorf("contract specifies strange range secret size: %v\n", pushes.SecretSize)
	}

	contractAddr, err := btcutil.NewAddressScriptHash(cmd.contract, coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return AuditedContract{}, err
	}

	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:], coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return AuditedContract{}, err
	}
	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return AuditedContract{}, err
	}

	contractValue := btcutil.Amount(cmd.contractTx.TxOut[contractOut].Value)

	contract := AuditedContract{Address: contractAddr, Value: contractValue, RecipientAddress: recipientAddr, RecipientRefundAddress: refundAddr, LockTime: pushes.LockTime, AtomicSwapDataPushes:pushes}

	return contract, nil
}

func (contract AuditedContract) show() error{
	fmt.Printf("Contract address:        %v\n", contract.Address)
	fmt.Printf("Contract value:          %v\n", contract.Value)
	fmt.Printf("Recipient address:       %v\n", contract.RecipientAddress)
	fmt.Printf("Recipient refund address: %v\n\n", contract.RecipientRefundAddress)

	if contract.AtomicSwapDataPushes.LockTime >= int64(txscript.LockTimeThreshold) {
		t := time.Unix(contract.AtomicSwapDataPushes.LockTime, 0)
		fmt.Printf("Locktime: %v\n", t.UTC())
		reachedAt := time.Until(t).Truncate(time.Second)
		if reachedAt > 0 {
			fmt.Printf("Locktime reached in %v\n", reachedAt)
		} else {
			fmt.Printf("Contract refund time lock has expired !\n")
			return nil
		}
		fmt.Printf("Locktime: block %v\n", contract.AtomicSwapDataPushes.LockTime)
	}
	return nil
}

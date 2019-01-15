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
	"strings"
	"time"
)

type AuditContractCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

type AuditedContract struct {
	Coin                   string         `json:"contract_coin"`
	Address                string         `json:"contract_address"`
	valueSat               btcutil.Amount `json:"contract_value_satoshi"`
	Value                  float64        `json:"contract_value"`
	ValueCoin              string         `json:"value_coin"`
	RecipientAddress       string         `json:"recipient_address"`
	AuthorRefundAddress string         `json:"author_refund_address"`
	SecretHash             string         `json:"secret_hash"`
	LockTime               int64          `json:"lock_time"`
	LockTimeReachedIn      time.Duration  `json:"lock_time_reached_in"`
	LockTimeExpired        bool           `json:"lock_time_expired"`
	atomicSwapDataPushes *txscript.AtomicSwapDataPushes
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
	cryptocurrencyName := coin.Name
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

	AuthorRefundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return AuditedContract{}, err
	}

	contractValue := btcutil.Amount(cmd.contractTx.TxOut[contractOut].Value)
	ValueCoin := fmt.Sprintf("%f %s", contractValue.ToBTC(), strings.ToUpper(coin.Symbol))

	var lockTimeExpired bool
	var lockTimeReachedIn time.Duration

	if pushes.LockTime >= int64(txscript.LockTimeThreshold) {
		t := time.Unix(pushes.LockTime, 0)
		fmt.Printf("Locktime: %v\n", t.UTC())
		reachedAt := time.Until(t).Truncate(time.Second)
		if reachedAt > 0 {
			lockTimeReachedIn = reachedAt
			fmt.Printf("Locktime reached in %v\n", reachedAt)
		} else {
			lockTimeExpired = true
			fmt.Printf("Contract refund time lock has expired !\n")
		}
		fmt.Printf("Locktime: block %v\n", pushes.LockTime)
	}

	contract := AuditedContract{
		Coin:                   cryptocurrencyName,
		Address:                contractAddr.EncodeAddress(),
		valueSat:               contractValue,
		Value:                  contractValue.ToBTC(),
		ValueCoin:              ValueCoin,
		RecipientAddress:       recipientAddr.EncodeAddress(),
		AuthorRefundAddress: AuthorRefundAddr.EncodeAddress(),
		SecretHash:             fmt.Sprintf("%x", pushes.SecretHash),
		LockTime:               pushes.LockTime,
		LockTimeReachedIn:      lockTimeReachedIn,
		LockTimeExpired:        lockTimeExpired,
		atomicSwapDataPushes:   pushes,
	}

	return contract, nil
}

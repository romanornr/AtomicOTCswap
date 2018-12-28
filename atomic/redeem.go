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
	"github.com/viacoin/viawallet/wallet/txrules"
	"strings"
)

type redeemCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
	secret     []byte
}

// coinTicker should be the coin the participant wants to redeem from the counter party
func Redeem(coinTicker string, contractHex string, contractTransaction string, secretHex string, wif *btcutil.WIF) error {
	coin, err := bcoins.SelectCoin(coinTicker)
	if err != nil {
		return err
	}

	contract, err := hex.DecodeString(contractHex)
	if err != nil {
		return fmt.Errorf("failed to decode contract: %v\n", err)
	}

	contractTxBytes, err := hex.DecodeString(contractTransaction)
	if err != nil {
		return fmt.Errorf("failed to decode contract transaction: %v\n", err)
	}
	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return fmt.Errorf("failed to decode contract transaction: %v\n", err)
	}

	secret, err := hex.DecodeString(secretHex)
	if err != nil {
		return fmt.Errorf("failed to decode secret: %v\n", err)
	}

	cmd := &redeemCmd{contract: contract, contractTx: &contractTx, secret: secret}
	return cmd.runRedeem(wif, coin)
}

func (cmd *redeemCmd) runRedeem(wif *btcutil.WIF, coin bcoins.Coin) error {
	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, cmd.contract)
	if err != nil {
		return err
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool")
	}
	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:], coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return err
	}
	contractHash := btcutil.Hash160(cmd.contract)
	contractOut := -1
	for i, out := range cmd.contractTx.TxOut {
		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, coin.Network.ChainCgfMainNetParams())
		if sc == txscript.ScriptHashTy &&
			bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash) {
			contractOut = i
			break
		}
	}

	if contractOut == -1 {
		return errors.New("transaction does not contain  a contract output")
	}

	addr, err := getRawChangeAddress(wif)
	if err != nil {
		return fmt.Errorf("getrawchangeAddress: %v\n", err)
	}
	outScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return err
	}

	contractTxHash := cmd.contractTx.TxHash()
	contractOutPoint := wire.OutPoint{
		Hash:  contractTxHash,
		Index: uint32(contractOut),
	}

	feePerKb, minFeePerKb, err := GetFeePerKB()
	if err != nil {
		return err
	}

	redeemTx := wire.NewMsgTx(coin.TxVersion)
	redeemTx.LockTime = uint32(pushes.LockTime)
	redeemTx.AddTxIn(wire.NewTxIn(&contractOutPoint, nil, nil))
	redeemTx.AddTxOut(wire.NewTxOut(0, outScript))
	redeemSize := estimateRedeemSerializeSize(cmd.contract, redeemTx.TxOut)
	fee := txrules.FeeForSerializeSize(feePerKb, redeemSize)
	redeemTx.TxOut[0].Value = cmd.contractTx.TxOut[contractOut].Value - int64(fee)
	if txrules.IsDustOutput(redeemTx.TxOut[0], minFeePerKb) {
		return fmt.Errorf("redeem output value of %v %s is dust", btcutil.Amount(redeemTx.TxOut[0].Value).ToBTC(), strings.ToUpper(coin.Symbol))
	}

	redeemSig, redeemPubKey, err := createSig(redeemTx, 0, cmd.contract, recipientAddr, wif)
	if err != nil {
		return err
	}
	redeemScriptSig, err := redeemP2SHContract(cmd.contract, redeemSig, redeemPubKey, cmd.secret)
	if err != nil {
		return err
	}
	redeemTx.TxIn[0].SignatureScript = redeemScriptSig

	redeemTxHash := redeemTx.TxHash()
	redeemFeePerKb := calcFeePerKb(fee, redeemTx.SerializeSize())

	var buf bytes.Buffer
	buf.Grow(redeemTx.SerializeSize())
	redeemTx.Serialize(&buf)

	if verify {
		e, err := txscript.NewEngine(cmd.contractTx.TxOut[contractOutPoint.Index].PkScript,
			redeemTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(redeemTx), cmd.contractTx.TxOut[contractOut].Value)
		if err != nil {
			panic(err)
		}
		err = e.Execute()
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Redeem fee: %v (%0.8f %s/kB)\n\n", fee, redeemFeePerKb, coin.Symbol)
	fmt.Printf("Redeem transaction (%v):\n", &redeemTxHash)
	fmt.Printf("%x\n\n", buf.Bytes())

	return nil
}

// estimateRedeemSerializeSize returns a worst case serialize size estimates for
// a transaction that redeems an atomic swap P2SH output.
func estimateRedeemSerializeSize(contract []byte, txOuts []*wire.TxOut) int {
	contractPush, err := txscript.NewScriptBuilder().AddData(contract).Script()
	if err != nil {
		panic(err)
	}
	contractPushSize := len(contractPush)

	// 12 additional bytes are for version, locktime & expiry
	return 12 + wire.VarIntSerializeSize(1) + wire.VarIntSerializeSize(uint64(len(txOuts))) +
		inputSize(redeemAtomicSwapSigScriptSize+contractPushSize) +
		sumOutputSerializeSizes(txOuts)
}

// redeemP2SHContract returns the signature script to redeem a contract output
// using the redeemer's signature and the initiator's secret.  This function
// assumes P2SH and appends the contract as the final data push.
func redeemP2SHContract(contract, sig, pubKey, secret []byte) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddData(sig)
	builder.AddData(pubKey)
	builder.AddData(secret)
	builder.AddInt64(1)
	builder.AddData(contract)
	return builder.Script()
}

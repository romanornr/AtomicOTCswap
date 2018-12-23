package atomic

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viad/txscript"
	"github.com/viacoin/viad/wire"
	btcutil "github.com/viacoin/viautil"
	"github.com/viacoin/viawallet/wallet/txrules"
)

func Refund(contractHex string, contractTransaction string) (*refundCmd, error) {
	contract, err := hex.DecodeString(contractHex)
	if err != nil {
		return &refundCmd{}, fmt.Errorf("failed to decode contract: %v\n", err)
	}
	contractTxBytes, err := hex.DecodeString(contractTransaction)
	if err != nil {
		return &refundCmd{}, fmt.Errorf("failed to decode contract transaction: %v\n", err)
	}
	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return &refundCmd{}, fmt.Errorf("failed to decode transaction: %v\n", err)
	}

	return &refundCmd{contract: contract, contractTx: &contractTx}, nil
}

func (cmd *refundCmd) runRefund() error {
	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, cmd.contract)
	if err != nil {
		return err
	}
	if pushes == nil {
		return error.New("contract is not an atomic swap script recognized bu this tool")
	}

	feePerKb, minFeePerKb, err := GetFeePerKB()
	if err != nil {
		return err
	}

	refundTx, refundFee, err := buildRefund(cmd.contract, cmd.contractTx, feePerKb, minFeePerKb)
	if err != nil {
		return err
	}

	//TODO
}

func buildRefund(contract []byte, contractTx *wire.MsgTx, feePerKb, minFeePerKb btcutil.Amount, refundAddress *btcutil.Address, wif *btcutil.WIF) (refundTx *wire.MsgTx, refundFee btcutil.Amount, err error)  {
	contractP2SH, err := btcutil.NewAddressScriptHash(contract, &chaincfg.MainNetParams)
	if err != nil {
		return nil, 0, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, 0, err
	}

	contractTxHash := contractTx.TxHash()
	contractOutpoint := wire.OutPoint{Hash:contractTxHash, Index: ^uint32(0)}
	for i, out := range contractTx.TxOut {
		if bytes.Equal(out.PkScript, contractP2SHPkScript) {
			contractOutpoint.Index = uint32(i)
			break
		}
	}
	if contractOutpoint.Index == ^uint32(0) {
		return nil, 0, errors.New("contract tx does not contain a P2SH contract payment")
	}

	//addr, err := btcutil.DecodeAddress(refundAddress, &chaincfg.MainNetParams)
	//if err != nil {
	//	return fmt.Errorf("error decoding refund address: %s\n", addr.String())
	//}
	refundOutScript, err := txscript.PayToAddrScript(*refundAddress)
	if err != nil {
		return err
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, refundOutScript)
	if err != nil {
		panic(err) /// TODO Something else? Idk yet
	}

	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], &chaincfg.MainNetParams)
	if err != nil {
		return nil, 0, err
	}

	refundTx = wire.NewMsgTx(txVersion) // TODO refactor: This can be different per coin. Viacoin is 2 but Zcoin is 1
	refundTx.LockTime = uint32(pushes.LockTime)
	refundTx.AddTxOut(wire.NewTxOut(0, refundOutScript))
	refundSize := estimateRefundSerializeSize(contract, refundTx.TxOut)
	refundFee = txrules.FeeForSerializeSize(feePerKb, refundSize)
	refundTx.TxOut[0].Value = contractTx.TxOut[contractOutpoint.Index].Value - int64(refundFee)
	if txrules.IsDustOutput(refundTx.TxOut[0], minFeePerKb) {
		return nil, 0, fmt.Errorf("refund output value of %v is dust", btcutil.Amount(refundTx.TxOut[0].Value))
	}

	txIn := wire.NewTxIn(&contractOutpoint, nil, nil)
	txIn.Sequence = 0
	refundTx.AddTxIn(txIn)

	refundSig, refundPubkey, err := createSig(refundTx, 0, contract, wif) //TODO signing
	if err != nil {
		return nil, 0, err
	}

	refundSigScript, err := refundP2SHContract(contract, refundSig, refundPubkey)
	if err != nil {
		return nil, 0, err
	}

	refundTx.TxIn[0].SignatureScript = refundSigScript

}

func createSig(tx *wire.MsgTx, idx int, pkScript []byte, wif *btcutil.WIF) (sig, pubkey []byte, err error) {
	sourceAddress, _  := GenerateNewPublicKey(*wif)
	if !bytes.Equal(sourceAddress.ScriptAddress(), pkScript) {
		return nil, nil, fmt.Errorf("error signing address: %s\n", sourceAddress)
	}

	sig, err = txscript.SignatureScript(tx, idx, pkScript, txscript.SigHashAll, wif.PrivKey, true)
	if err != nil {
		return nil, nil, err
	}
	return sig, wif.PrivKey.PubKey().SerializeCompressed(), nil
}


// TODO---------------------------

const (
	redeemAtomicSwapSigScriptSize = 1 + 73 + 1 + 33 + 1 + 32 + 1
	refundAtomicSwapSigScriptSize = 1 + 73 + 1 + 33 + 1
)

func sumOutputSerializeSizes(outputs []*wire.TxOut) (serializeSize int) {
	for _, txOut := range outputs {
		serializeSize += txOut.SerializeSize()
	}
	return serializeSize
}


// inputSize returns the size of the transaction input needed to include a
// signature script with size sigScriptSize.  It is calculated as:
//
//   - 32 bytes previous tx
//   - 4 bytes output index
//   - Compact int encoding sigScriptSize
//   - sigScriptSize bytes signature script
//   - 4 bytes sequence
func inputSize(sigScriptSize int) int {
	return 32 + 4 + wire.VarIntSerializeSize(uint64(sigScriptSize)) + sigScriptSize + 4
}

// estimateRefundSerializeSize returns a worst case serialize size estimates for
// a transaction that refunds an atomic swap P2SH output.
func estimateRefundSerializeSize(contract []byte, txOuts []*wire.TxOut) int {
	contractPush, err := txscript.NewScriptBuilder().AddData(contract).Script()
	if err != nil {
		// Should never be hit since this script does exceed the limits.
		panic(err)
	}
	contractPushSize := len(contractPush)

	// 12 additional bytes are for version, locktime and expiry.
	return 12 + wire.VarIntSerializeSize(1) +
		wire.VarIntSerializeSize(uint64(len(txOuts))) +
		inputSize(refundAtomicSwapSigScriptSize+contractPushSize) +
		sumOutputSerializeSizes(txOuts)
}
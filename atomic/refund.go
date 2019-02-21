package atomic

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/romanornr/AtomicOTCswap/swaputil"
	"github.com/viacoin/viad/txscript"
	"github.com/viacoin/viad/wire"
	btcutil "github.com/viacoin/viautil"
	"github.com/viacoin/viawallet/wallet/txrules"
)

type refundCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

func Refund(contractHex string, contractTransaction string, wif *btcutil.WIF) (*refundCmd, error) {
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

func (cmd *refundCmd) Run(wif *btcutil.WIF, coin *bcoins.Coin) error {
	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, cmd.contract)
	if err != nil {
		return err
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized bu this tool")
	}

	feePerKb, minFeePerKb, err := GetFeePerKB()
	if err != nil {
		return err
	}

	refundTx, refundFee, err := buildRefund(cmd.contract, cmd.contractTx, feePerKb, minFeePerKb, wif, coin)
	if err != nil {
		return err
	}

	refundTxHash := refundTx.TxHash()
	var buf bytes.Buffer
	buf.Grow(refundTx.SerializeSize())
	refundTx.Serialize(&buf)

	refundFeePerKb := calcFeePerKb(refundFee, refundTx.SerializeSize())

	fmt.Printf("Refund fee: %v (%0.8f VIA/kB\n\n", refundFee, refundFeePerKb)
	fmt.Printf("Refund Transaction: (%v):\n", &refundTxHash)

	return nil
}

func calcFeePerKb(absoluteFee btcutil.Amount, serializeSize int) float64 {
	return float64(absoluteFee) / float64(serializeSize) / 1e5
}

func buildRefund(contract []byte, contractTx *wire.MsgTx, feePerKb, minFeePerKb btcutil.Amount, wif *btcutil.WIF, coin *bcoins.Coin) (refundTx *wire.MsgTx, refundFee btcutil.Amount, err error) {
	contractP2SH, err := btcutil.NewAddressScriptHash(contract, coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return nil, 0, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, 0, err
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{Hash: contractTxHash, Index: ^uint32(0)}
	for i, out := range contractTx.TxOut {
		if bytes.Equal(out.PkScript, contractP2SHPkScript) {
			contractOutPoint.Index = uint32(i)
			break
		}
	}
	if contractOutPoint.Index == ^uint32(0) {
		return nil, 0, errors.New("contract tx does not contain a P2SH contract payment")
	}

	address, err := getRawChangeAddress(wif, coin)
	if err != nil {
		fmt.Println(err)
	}

	refundOutScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return nil, 0, err
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		//expected only to be called with good input
		panic(err)
	}

	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return nil, 0, err
	}

	refundTx = wire.NewMsgTx(coin.TxVersion)
	refundTx.LockTime = uint32(pushes.LockTime)
	refundTx.AddTxOut(wire.NewTxOut(0, refundOutScript)) // amount set below
	refundSize := estimateRefundSerializeSize(contract, refundTx.TxOut)
	refundFee = txrules.FeeForSerializeSize(feePerKb, refundSize)
	refundTx.TxOut[0].Value = contractTx.TxOut[contractOutPoint.Index].Value - int64(refundFee)
	if txrules.IsDustOutput(refundTx.TxOut[0], minFeePerKb) {
		return nil, 0, fmt.Errorf("refund output value of %v is dust", btcutil.Amount(refundTx.TxOut[0].Value))
	}

	txIn := wire.NewTxIn(&contractOutPoint, nil, nil)
	txIn.Sequence = 0
	refundTx.AddTxIn(txIn)

	refundSig, refundPubKey, err := createSig(refundTx, 0, contract, refundAddr, wif, coin) //TODO signing
	if err != nil {
		return nil, 0, err
	}

	refundSigScript, err := refundP2SHContract(contract, refundSig, refundPubKey)
	if err != nil {
		return nil, 0, err
	}

	refundTx.TxIn[0].SignatureScript = refundSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			refundTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(refundTx), contractTx.TxOut[contractOutPoint.Index].Value)
		if err != nil {
			panic(err)
		}
		err = e.Execute()
		if err != nil {
			panic(err)
		}
	}
	return refundTx, refundFee, nil
}

// refundP2SHContract returns the signature script to refund a contract output
// using the contract initiator/participant signature after locktime is reached
func refundP2SHContract(contract, sig, pubkey []byte) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddData(sig)
	builder.AddData(pubkey)
	builder.AddInt64(0)
	builder.AddData(contract)
	return builder.Script()
}

func createSig(tx *wire.MsgTx, idx int, pkScript []byte, addr btcutil.Address, wif *btcutil.WIF, coin *bcoins.Coin) (sig, pubkey []byte, err error) {
	sourceAddress, _ := swaputil.GenerateNewPublicKey(*wif, coin.Network.ChainCgfMainNetParams())
	if sourceAddress.EncodeAddress() != addr.EncodeAddress() {
		return nil, nil, fmt.Errorf("error signing address: %s\n", sourceAddress)
	}

	sig, err = txscript.RawTxInSignature(tx, idx, pkScript, txscript.SigHashAll, wif.PrivKey)
	if err != nil {
		return nil, nil, err
	}
	return sig, wif.PrivKey.PubKey().SerializeCompressed(), nil
}

// getRawChangeAddress calls the getrawchangeaddress JSON-RPC method.  It is
// implemented manually as the rpcclient implementation always passes the
// account parameter which was removed in Viacoin Core 0.15.
func getRawChangeAddress(wif *btcutil.WIF, coin *bcoins.Coin) (btcutil.Address, error) {

	addr, _ := swaputil.GenerateNewPublicKey(*wif, coin.Network.ChainCgfMainNetParams())
	address, err := btcutil.DecodeAddress(addr.EncodeAddress(), coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return nil, err
	}

	if _, ok := address.(*btcutil.AddressPubKeyHash); !ok {
		return nil, fmt.Errorf("getrawchangeaddress: address %v is not P2PKH",
			addr)
	}
	return address, nil
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

package atomic

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/golangcrypto/ripemd160"
	"github.com/go-errors/errors"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	"github.com/romanornr/AtomicOTCswap/insight"
	"github.com/viacoin/viad/chaincfg/chainhash"
	"github.com/viacoin/viad/txscript"
	"github.com/viacoin/viad/wire"
	btcutil "github.com/viacoin/viautil"
)

const (
	verify     = true
	secretSize = 32
	txVersion  = 2
)

type Command struct {
	Command string
	Params  []string
}

type contractArgs struct {
	coin       *bcoins.Coin
	them       *btcutil.AddressPubKeyHash
	amount     btcutil.Amount
	locktime   int64
	secretHash []byte
}

// builtContract houses the details regarding a contract and the contract
// payment transaction, as well as the transaction to perform a refund.
type builtContract struct {
	contract       []byte
	contractP2SH   btcutil.Address
	contractTxHash *chainhash.Hash
	contractTx     *wire.MsgTx
	contractFee    btcutil.Amount
	refundTx       *wire.MsgTx
	refundFee      btcutil.Amount
}

func buildContract(args *contractArgs, wif *btcutil.WIF) (*builtContract, error) {

	refundAddress, _ := GenerateNewPublicKey(*wif, args.coin)
	refundAddr, _ := btcutil.DecodeAddress(refundAddress.EncodeAddress(), args.coin.Network.ChainCgfMainNetParams())

	refundAddrHash, ok := refundAddr.(interface {
		Hash160() *[ripemd160.Size]byte
	})

	if !ok {
		return nil, errors.New("unable to create hash160 from change address")
	}
	contract, err := atomicSwapContract(refundAddrHash.Hash160(), args.them.Hash160(), args.locktime, args.secretHash)
	if err != nil {
		return nil, err
	}
	contractP2SH, err := btcutil.NewAddressScriptHash(contract, args.coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return nil, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, err
	}

	feePerKb, minFeePerKb, err := GetFeePerKB()
	if err != nil {
		return nil, err
	}

	unsignedContract := wire.NewMsgTx(txVersion)
	unsignedContract.AddTxOut(wire.NewTxOut(int64(args.amount), contractP2SHPkScript))

	contractTx, contractFee, complete, err := fundAndSignRawTransaction(unsignedContract, wif, args.amount, args.coin)
	if err != nil {
		return nil, fmt.Errorf("signrawtransaction: %v", err)
	}
	if !complete {
		return nil, errors.New("signrawtransaction: failed to completely sign contract transaction")
	}

	contractTxHash := contractTx.TxHash()

	refundTx, refundFee, err := buildRefund(contract, contractTx, feePerKb, minFeePerKb, wif, args.coin)
	if err != nil {
		return nil, err
	}

	return &builtContract{
		contract,
		contractP2SH,
		&contractTxHash,
		contractTx,
		contractFee,
		refundTx,
		refundFee,
	}, nil
}

func fundAndSignRawTransaction(tx *wire.MsgTx, wif *btcutil.WIF, amount btcutil.Amount, coin *bcoins.Coin) (*wire.MsgTx, btcutil.Amount, bool, error) {
	sourceAddress, _ := GenerateNewPublicKey(*wif, coin)
	sourcePKScript, err := txscript.PayToAddrScript(sourceAddress.AddressPubKeyHash())
	if err != nil {
		return tx, amount, false, err
	}

	unspentOutputs := insight.GetUnspentOutputs(sourceAddress.AddressPubKeyHash().String(), coin)
	sourceUTXOs := insight.GetMinimalRequiredUTXO(int64(amount), unspentOutputs)
	var availableAmountToSpend int64

	// combine inputs
	for idx := range sourceUTXOs {
		availableAmountToSpend += sourceUTXOs[idx].Amount
		sourceUTXO := wire.NewOutPoint(sourceUTXOs[idx].Hash, sourceUTXOs[idx].TxIndex)
		sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)
		tx.AddTxIn(sourceTxIn)
	}

	change := availableAmountToSpend - int64(amount)
	var fee btcutil.Amount

	if change < 0 {
		return &wire.MsgTx{}, 0, false, fmt.Errorf("not enough funds to spend, Available amount %f %s", btcutil.Amount(availableAmountToSpend).ToBTC(), coin.Unit)
	}


		changeAddress := sourceAddress
		changeSendToScript, err := txscript.PayToAddrScript(changeAddress)
		if err != nil {
			return &wire.MsgTx{}, 0, false, fmt.Errorf("change address wrong\n")
		}

	if change >= 0 {
		changeOutput := wire.NewTxOut(change, changeSendToScript)
		//if change < 0 {
		//	maxAmountAvailable := btcutil.Amount(change) + amount
		//	return &wire.MsgTx{}, 0, false, fmt.Errorf("not enough funds to cover the fee of %f %s. Try %v %s", fee.ToBTC(), coin.Unit, maxAmountAvailable.ToBTC(), coin.Unit)
		//}
		changeOutput = wire.NewTxOut(change, changeSendToScript)
		tx.AddTxOut(changeOutput) // TODO NO CHANGE OUTPUT IF ITS DUST
	}

	// sign all transactions
	for i := range sourceUTXOs {
		sigScript, err := txscript.SignatureScript(tx, i, sourcePKScript, txscript.SigHashAll, wif.PrivKey, true)
		if err != nil {
			fmt.Errorf("error signing source UTXO's\n")
		}
		tx.TxIn[i].SignatureScript = sigScript
	}

	fee = feeEstimationBySize(tx.SerializeSize(), coin)
	change -= int64(fee)
	tx.TxOut[1] = wire.NewTxOut(change, changeSendToScript)

	for i := range sourceUTXOs {
		sigScript, err := txscript.SignatureScript(tx, i, sourcePKScript, txscript.SigHashAll, wif.PrivKey, true)
		if err != nil {
			fmt.Errorf("error signing source UTXO's\n")
		}
		tx.TxIn[i].SignatureScript = sigScript
	}

	signedTx := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(signedTx); err != nil {
		return &wire.MsgTx{}, 0, false, fmt.Errorf("failed to sign tx")
	}
	fmt.Printf("size: %d", tx.SerializeSize())

	return tx, fee, true, nil
}

//// TODO maybe fee per byte to 279
//func feeEstimator(tx *wire.MsgTx, ) (amount btcutil.Amount) {
//	feePerByte := 110 // TODO change for alts
//	estimatedSize := tx.SerializeSize()
//	return btcutil.Amount(feePerByte * estimatedSize)
//}

func feeEstimationBySize(size int, coin *bcoins.Coin) (amount btcutil.Amount) {
	feePerByte := coin.FeePerByte
	return btcutil.Amount(feePerByte * size)
}

// atomicSwapContract returns an output script that may be redeemed by one of 2 signature scripts:
// <their sig> <their pubkey> <initiator secret> 1
// <my sig> <my pubkey> 0
func atomicSwapContract(pkhMe, pkhThem *[ripemd160.Size]byte, locktime int64, secretHash []byte) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_IF) // if top of stack value is not False, execute. The top stack value is removed.
	{
		// require initiator's secret to be a known length that the redeeming party can audit.
		// this is used to prevent fraud attacks between 2 currencies that have different maximum data sizes
		builder.AddOp(txscript.OP_SIZE)        // pushes the string length of the top element of the stack (without popping it)
		builder.AddInt64(secretSize)           // pushes initiator secret length
		builder.AddOp(txscript.OP_EQUALVERIFY) // if inputs are equal, mark tx as valid

		// require initiator's secret to be known to redeem the output
		builder.AddOp(txscript.OP_SHA256)      // pushes the length of a SHA25 size
		builder.AddData(secretHash)            // push the data to the end of the script
		builder.AddOp(txscript.OP_EQUALVERIFY) // if inputs are equal, mark tx as valid

		// verify their signature is used to redeem the ouput
		// normally it ends with OP_EQUALVERIFY OP_CHECKSIG but
		// this has been moved outside of the branch to save a couple bytes
		builder.AddOp(txscript.OP_DUP)     // duplicates the stack of the top item
		builder.AddOp(txscript.OP_HASH160) // input has been hashed with SHA-256 and then with RIPEMD160 after
		builder.AddData(pkhThem[:])        // push the data to the end of the script
	}

	builder.AddOp(txscript.OP_ELSE) // refund path
	{
		// verify the locktime & drop if off the stack
		builder.AddInt64(locktime)                     // pushes locktime
		builder.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY) // verify locktime
		builder.AddOp(txscript.OP_DROP)                // remove the top stack item (locktime)

		// verify our signature is being used to redeem the output
		// normally it ends with OP_EQUALVERIFY OP_CHECKSIG but
		// this has been moved outside of the branch to save a couple bytes
		builder.AddOp(txscript.OP_DUP)     // duplicates the stack of the top item
		builder.AddOp(txscript.OP_HASH160) // input has been hashed with SHA-256 and then with RIPEMD160 after
		builder.AddData(pkhMe[:])          // push the data to the end of the script

	}
	builder.AddOp(txscript.OP_ENDIF) // all blocks must end, or the transaction is invalid

	// returns 1 if the inputs are exactly equal, 0 otherwise.
	// mark transaction as invalid if top of stack is not true. The top stack value is removed.
	builder.AddOp(txscript.OP_EQUALVERIFY)

	// The entire transaction's outputs, inputs, and script are hashed.
	// The signature used by OP_CHECKSIG must be a valid signature for this hash
	// and public key. If it is, 1 is returned, 0 otherwise.
	builder.AddOp(txscript.OP_CHECKSIG)
	return builder.Script()
}

func sha256Hash(x []byte) []byte {
	hash := sha256.Sum256(x)
	return hash[:]
}

func GetFeePerKB() (useFee, replayFee btcutil.Amount, err error) {

	// TODO fix for multiple coins

	relayFee, _ := btcutil.NewAmount(0.00100000) //rpc call -> getnetworkinfo: relayfee
	payTxFee, _ := btcutil.NewAmount(0.00000000) //rpc call -> getwalletinfo: paytxfee

	if payTxFee != 0 {
		maxFee := payTxFee
		if relayFee > maxFee {
			maxFee = relayFee
		}
		return maxFee, relayFee, nil
	}

	return relayFee, relayFee, nil
}

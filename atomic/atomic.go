package atomic

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/romanornr/AtomicOTCswap/insight"
	"time"

	"github.com/btcsuite/golangcrypto/ripemd160"
	"github.com/go-errors/errors"
	"github.com/viacoin/viad/chaincfg"
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

//type Command interface {
//
//}

type initiateCmd struct {
	counterparty2Addr *btcutil.AddressPubKeyHash
	refundAddr	*btcutil.AddressPubKeyHash
	wif *btcutil.WIF
	amount            btcutil.Amount
}

type participateCmd struct {
	counterparty1Addr *btcutil.AddressPubKeyHash
	amount            btcutil.Amount
}

type redeemCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

type refundCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

type extractSecretCmd struct {
	redemptionTx *wire.MsgTx
	secretHash   []byte
}

type Command struct {
	Command string
	Params  []string
}

type contractArgs struct {
	them       *btcutil.AddressPubKeyHash
	amount     btcutil.Amount
	locktime   int64
	secretHash []byte
}

func Initiate(participantAddr string, wif *btcutil.WIF, amount float64) error {
	counterparty2Addr, err := btcutil.DecodeAddress(participantAddr, &chaincfg.MainNetParams)
	if err != nil {
		return fmt.Errorf("failed to decode the address from the participant: %s", err)
	}

	counterparty2AddrP2KH, ok := counterparty2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return errors.New("participant address is not P2KH")
	}

	amount2, err := btcutil.NewAmount(amount)
	if err != nil {
		return err
	}

	refundAddrPubKey, err := GenerateNewPublicKey(*wif)

	refundAddr, err := btcutil.DecodeAddress(refundAddrPubKey.EncodeAddress(), &chaincfg.MainNetParams)
	if err != nil {
		return fmt.Errorf("failed to decode the refund address: %s", err)
	}


	refundAddrP2KH, ok := refundAddr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return errors.New("participant address is not P2KH")
	}

	cmd := &initiateCmd{counterparty2Addr: counterparty2AddrP2KH, refundAddr: refundAddrP2KH, wif: wif, amount: amount2}
	return cmd.runCommand()
}

func (cmd *initiateCmd) runCommand() error {
	var secret [secretSize]byte
	_, err := rand.Read(secret[:])
	if err != nil {
		return err
	}

	secretHash := sha256Hash(secret[:])
	locktime := time.Now().Add(5 * time.Minute).Unix() // NEED TO CHANGE

	build, err := buildContract(&contractArgs{
		them:       cmd.counterparty2Addr,
		amount:     cmd.amount,
		locktime:   locktime,
		secretHash: secretHash,
	}, cmd.wif)

	if err != nil {
		return err
	}

	refundTxHash := build.refundTx.TxHash()

	fmt.Printf("Secret:      %x\n", secret)
	fmt.Printf("Secret hash: %x\n\n", secretHash)
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

	return nil
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

	refundAddress, _ := GenerateNewPublicKey(*wif)
	refundAddr, _ := btcutil.DecodeAddress(refundAddress.EncodeAddress(), &chaincfg.MainNetParams)

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
	contractP2SH, err := btcutil.NewAddressScriptHash(contract, &chaincfg.MainNetParams)
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
	unsignedContract, contractFee, sourceUTXOs, err := fundRawTransaction(unsignedContract, wif, feePerKb)
	if err != nil {
		return nil, fmt.Errorf("funded raw transaction: %v\n", err)
	}

	contractTx, complete, err := signRawTransaction(unsignedContract, wif, sourceUTXOs)
	if err != nil {
		return nil, fmt.Errorf("signrawtransaction: %v", err)
	}
	if !complete {
		return nil, errors.New("signrawtransaction: failed to completely sign contract transaction")
	}

	contractTxHash := contractTx.TxHash()

	refundTx, refundFee, err := buildRefund(contract, contractTx, feePerKb, minFeePerKb, wif)
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

var AmountPaymentSat = int64(2.09 * 10000000)

func fundRawTransaction(tx *wire.MsgTx, wif *btcutil.WIF, feePerKb btcutil.Amount) (fundedTx *wire.MsgTx, fee btcutil.Amount, sourceUTXOs []*insight.UTXO, err error) {

	sourceAddress, _ := GenerateNewPublicKey(*wif)
	unspentOutputs := insight.GetUnspentOutputs(sourceAddress.AddressPubKeyHash().String())
	sourceUTXOs = insight.GetMinimalRequiredUTXO(AmountPaymentSat, unspentOutputs)
	availableAmountToSpend := int64(0) // amount in UTXO available

	for idx := range sourceUTXOs {
		availableAmountToSpend += sourceUTXOs[idx].Amount
		sourceUTXO := wire.NewOutPoint(sourceUTXOs[idx].Hash, sourceUTXOs[idx].TxIndex)
		sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)
		tx.AddTxIn(sourceTxIn)
	}

	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())
	fmt.Println(tx.SerializeSize())
	tx.Serialize(&buf)

	//feeAmount, err := btcutil.NewAmount(0.00045)
	//if err != nil {
	//	return nil, 0, sourceUTXOs, err
	//}

	feeAmount := feeEstimator(tx)

	fmt.Println(hex.EncodeToString(buf.Bytes()))
	return tx, feeAmount, sourceUTXOs, nil
}

func feeEstimator(tx *wire.MsgTx) (amount btcutil.Amount) {
	feePerByte := 110 // TODO change for alts
	estimatedSize := tx.SerializeSize()

	a := btcutil.Amount(feePerByte * estimatedSize)
	fmt.Println(a)
	return a
}

func signRawTransaction(tx *wire.MsgTx, wif *btcutil.WIF, sourceUTXOs []*insight.UTXO) (*wire.MsgTx, bool, error) {

	sourceAddress, _ := GenerateNewPublicKey(*wif)
	sourcePKScript, err := txscript.PayToAddrScript(sourceAddress.AddressPubKeyHash())
	if err != nil {
		fmt.Errorf("error signing soucePKScript: %s\n", sourcePKScript)
	}

	for i := range sourceUTXOs {
		sigScript, err := txscript.SignatureScript(tx, i, sourcePKScript, txscript.SigHashAll, wif.PrivKey, true)
		if err != nil {
			fmt.Errorf("error signing source UTXO's\n")
		}
		tx.TxIn[i].SignatureScript = sigScript
	}

	signedTx := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(signedTx); err != nil {
		return &wire.MsgTx{}, false, fmt.Errorf("Failed to sign tx")
	}

	fmt.Printf("signed tx: %s\n", hex.EncodeToString(signedTx.Bytes()))
	return tx, true, nil
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

	relayFee, _ := btcutil.NewAmount(0.00000001) //https://github.com/viacoin/viacoin/blob/master/src/test/amount_tests.cpp#L105
	payTxFee, _ :=  btcutil.NewAmount(0.00000000) //rpc call -> getwalletinfo: paytxfee

	if payTxFee != 0 {
		maxFee := payTxFee
		if relayFee > maxFee {
			maxFee = relayFee
		}
		return maxFee, relayFee, nil
	}

	return relayFee, relayFee, nil
}

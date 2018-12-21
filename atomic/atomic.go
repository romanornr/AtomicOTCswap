package atomic

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/golangcrypto/ripemd160"
	"github.com/go-errors/errors"
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viad/chaincfg/chainhash"
	"github.com/viacoin/viad/txscript"
	"github.com/viacoin/viad/wire"
	btcutil "github.com/viacoin/viautil"
	"time"
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

func Initiate(participantAddr string, refundAddr string, amount float64) error {
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

	refundAddress, err := btcutil.DecodeAddress(refundAddr, &chaincfg.MainNetParams)
	if err != nil {
		return fmt.Errorf("failed to decode the refund address: %s", err)
	}

	cmd := &initiateCmd{counterparty2Addr: counterparty2AddrP2KH, amount: amount2}
	return cmd.runCommand(refundAddress)
}

func (cmd *initiateCmd) runCommand(refundAddr btcutil.Address) error {
	var secret [secretSize]byte
	_, err := rand.Read(secret[:])
	if err != nil {
		return err
	}

	secretHash := sha256Hash(secret[:])
	locktime := time.Now().Add(10 * time.Minute).Unix() // NEED TO CHANGE TO 48 HOURS


	build, err := buildContract(&contractArgs {
		them:       cmd.counterparty2Addr,
		amount:     cmd.amount,
		locktime:   locktime,
		secretHash: secretHash,
	}, refundAddr)

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



func buildContract(args *contractArgs, refundAddr btcutil.Address) (*builtContract, error) {

	//refundAddr, _ := btcutil.DecodeAddress("VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A", &chaincfg.MainNetParams)

	refundAddrHash, ok := refundAddr.(interface {
		Hash160() *[ripemd160.Size]byte
	})

	if !ok {
		return nil, errors.New("unable to create hash160 from change address")
	}
	contract, err  := atomicSwapContract(refundAddrHash.Hash160(), args.them.Hash160(), args.locktime, args.secretHash)
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
	feePerKb, _ := btcutil.NewAmount(0.01)
    //minFeePerKb, _ := btcutil.NewAmount(0.02)

	unsignedContract := wire.NewMsgTx(txVersion)
	unsignedContract.AddTxOut(wire.NewTxOut(int64(args.amount), contractP2SHPkScript))
	unsignedContract, contractFee, err := fundRawTransaction(unsignedContract, feePerKb)
	if err != nil {
		return nil, fmt.Errorf("funded raw transaction: %v\n", err)
	}
	//contractTx, complete, err := signRawTransaction(unsignedContract)
	contractTx, complete, err := signRawTransaction(unsignedContract)
	if err != nil {
		return nil, fmt.Errorf("signrawtransaction: %v", err)
	}
	if !complete {
		return nil, errors.New("signrawtransaction: failed to completely sign contract transaction")
	}

	contractTxHash := contractTx.TxHash()

	//refundTx, refundFee, err := buildRefund(c, contract, contractTx, feePerKb, minFeePerKb)
	//if err != nil {
	//	return nil, err
	//}

	refundFee, _  := btcutil.NewAmount(0.001)


	return &builtContract{
		contract,
		contractP2SH,
		&contractTxHash,
		contractTx,
		contractFee,
		&wire.MsgTx{},
		//refundTx,
		refundFee,
	}, nil
}

func fundRawTransaction(tx *wire.MsgTx, feePerKb btcutil.Amount) (fundedTx *wire.MsgTx, fee btcutil.Amount, err error) {
	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())
	tx.Serialize(&buf)

	//param0 := hex.EncodeToString(buf.Bytes())
	//param0, err := json.Marshal(hex.EncodeToString(buf.Bytes()))
	//if err != nil {
	//	return nil, 0, err
	//}

	//param1, err := json.Marshal(struct {
	//	FeeRate float64 `json:"feeRate"`
	//}{
	//	FeeRate: feePerKb.ToBTC(),
	//})

	var funded struct{
		Hex string `json:"hex"`
		Fee float64 `json:"fee"`
		ChangePos float64 `json:"chane_pos"`
	}
	//x := hex.EncodeToString(buf.Bytes())
	funded.Fee = 0.001
	//funded.Hex = x
	//
	//fmt.Println(buf.Bytes())

	///fundedTxBytes, err := hex.DecodeString(param0) // TODO
	//fundedTx = &wire.MsgTx{}
	//err = fundedTx.Deserialize(bytes.NewReader(fundedTxBytes))
	//if err != nil {
	//	return nil, 0, err /// TODO fix
	//}

	feeAmount, err := btcutil.NewAmount(funded.Fee)
	if err != nil {
		return nil, 0, err
	}

	fmt.Println(hex.EncodeToString(buf.Bytes()))
	return tx, feeAmount, nil
}

func signRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	signedTx := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(signedTx); err != nil {
		return &wire.MsgTx{}, false, fmt.Errorf("Failed to sign tx")
	}

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

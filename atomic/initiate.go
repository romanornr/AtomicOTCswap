package atomic

import (
	"bytes"
	"fmt"
	"github.com/viacoin/viad/chaincfg"
	btcutil "github.com/viacoin/viautil"
	"time"
)
type initiateCmd struct {
	counterparty2Addr *btcutil.AddressPubKeyHash
	//refundAddr	*btcutil.AddressPubKeyHash
	//wif *btcutil.WIF
	amount            btcutil.Amount
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

	//refundAddrPubKey, err := GenerateNewPublicKey(*wif)
	//
	//refundAddr, err := btcutil.DecodeAddress(refundAddrPubKey.EncodeAddress(), &chaincfg.MainNetParams)
	//if err != nil {
	//	return fmt.Errorf("failed to decode the refund address: %s", err)
	//}
	//
	//
	//refundAddrP2KH, ok := refundAddr.(*btcutil.AddressPubKeyHash)
	//if !ok {
	//	return errors.New("participant address is not P2KH")
	//}

	cmd := &initiateCmd{counterparty2Addr: counterparty2AddrP2KH, amount: amount2}
	return cmd.runCommand(wif)
}

func (cmd *initiateCmd) runCommand(wif *btcutil.WIF) error {
	var secret [secretSize]byte
	_, err := rand.Read(secret[:])
	if err != nil {
		return err
	}

	secretHash := sha256Hash(secret[:])
	locktime := time.Now().Add(10 * time.Minute).Unix() // NEED TO CHANGE

	build, err := buildContract(&contractArgs{
		them:       cmd.counterparty2Addr,
		amount:     cmd.amount,
		locktime:   locktime,
		secretHash: secretHash,
	}, wif)

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

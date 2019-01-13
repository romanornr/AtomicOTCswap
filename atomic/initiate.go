package atomic

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/romanornr/AtomicOTCswap/bcoins"
	btcutil "github.com/viacoin/viautil"
	"strings"
	"time"
)
type initiateCmd struct {
	counterParty2Addr *btcutil.AddressPubKeyHash
	amount            btcutil.Amount
}


func Initiate(coinTicker string, participantAddr string, wif *btcutil.WIF, amount float64) error {

	coin, err := bcoins.SelectCoin(coinTicker)
	if err != nil {
		return err
	}

	fmt.Println(coin.Network.ChainCgfMainNetParams().PubKeyHashAddrID)
	fmt.Println(coin.Network.P2PKH)

	counterParty2Addr, err := btcutil.DecodeAddress(participantAddr, coin.Network.ChainCgfMainNetParams())
	if err != nil {
		return fmt.Errorf("failed to decode the address from the participant: %s\n", err)
	}

	fmt.Println("gi")

	counterParty2AddrP2KH, ok := counterParty2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return errors.New("participant address is not P2KH")
	}

	amount2, err := btcutil.NewAmount(amount)
	if err != nil {
		return err
	}

	cmd := &initiateCmd{counterParty2Addr: counterParty2AddrP2KH, amount: amount2}
	return cmd.runCommand(wif, &coin, amount)
}

func (cmd *initiateCmd) runCommand(wif *btcutil.WIF, coin *bcoins.Coin, amount float64) error {
	var secret [secretSize]byte
	_, err := rand.Read(secret[:])
	if err != nil {
		return err
	}

	secretHash := sha256Hash(secret[:])
	locktime := time.Now().Add(10 * time.Minute).Unix() // NEED TO CHANGE

	build, err := buildContract(&contractArgs{
		coin1: coin,
		them:       cmd.counterParty2Addr,
		amount:     cmd.amount,
		locktime:   locktime,
		secretHash: secretHash,
	}, wif)

	if err != nil {
		return err
	}

	ticker := strings.ToUpper(coin.Symbol)
	refundTxHash := build.refundTx.TxHash()

	fmt.Printf("Secret:      %x\n", secret)
	fmt.Printf("Secret hash: %x\n\n", secretHash)

	fmt.Printf("Contract amount: %v %s\n", amount, ticker)
	fmt.Printf("Contract fee: %v %s\n", build.contractFee.ToBTC(), ticker)
	fmt.Printf("Refund fee:   %v %s\n", build.refundFee.ToBTC(), ticker)
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

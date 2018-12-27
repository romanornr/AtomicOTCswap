package atomic

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/viacoin/viad/wire"
	btcutil "github.com/viacoin/viautil"
)

type redeemCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
	secret     []byte
}

func Redeem(contractHex string, contractTransaction string, secretHex string, wif *btcutil.WIF) error {
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
	return cmd.runRedeem(wif)
}

func (cmd *redeemCmd) runRedeem(wif *btcutil.WIF) error { // TODO
	//pushes, err := txscript.ExtractAtomicSwapDataPushes(0, cmd.contract)
	//if err != nil {
	//	return err
	//}
	//if pushes == nil {
	//	return errors.New("contract is not an atomic swap script recognized by this tool")
	//}
	//recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:], &chaincfg.MainNetParams)
	//if err != nil {
	//	return err
	//}
	//contractHash := btcutil.Hash160(cmd.contract)
	//contractOut := -1
	//for i, out := range cmd.contractTx.TxOut {
	//	sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.MainNetParams)
	//	if sc == txscript.ScriptHashTy &&
	//		bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash) {
	//		contractOut = i
	//		break
	//	}
	//}
	//
	//if contractOut == -1 {
	//	return errors.New("transaction does not contain  a contract output")
	//}
	//
	//addr, err := getRawChangeAddress(wif)
	//if err != nil {
	//	return fmt.Errorf("getrawchangeAddress: %v\n", err)
	//}
	//outScript, err := txscript.PayToAddrScript(addr)
	//if err != nil {
	//	return err
	//}
	//
	//contractTxHash := cmd.contractTx.TxHash()
	//contractOutPoint := wire.OutPoint{
	//	Hash:contractTxHash,
	//	Index:uint32(contractOut),
	//}
	//
	//feePerKb, minFeePerKb, err := GetFeePerKB()
	//if err != nil {
	//	return err
	//}
	//

	return nil
}

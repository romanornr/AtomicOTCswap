package atomic

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/viacoin/viad/txscript"
	"github.com/viacoin/viad/wire"
)

type extractSecretCmd struct {
	redemptionTx *wire.MsgTx
	secretHash   []byte
}

func extractSecret(redemptionTransaction string, secret string) (string, error) {
	redemptionTxBytes, err := hex.DecodeString(redemptionTransaction)
	if err != nil {
		return "", fmt.Errorf("failed to decode redemption transaction: %v\n", err)
	}

	var redemptionTx wire.MsgTx
	err = redemptionTx.Deserialize(bytes.NewReader(redemptionTxBytes))
	if err != nil {
		return "", fmt.Errorf("failed to decode redemptioon transaction: %v\n", err)
	}

	secretHash, err := hex.DecodeString(secret)
	if err != nil {
		return "", errors.New("secret hash must be hex encoded")
	}
	if len(secretHash) != sha256.Size {
		return "", errors.New("secret hash wrong size")
	}

	cmd := &extractSecretCmd{redemptionTx: &redemptionTx, secretHash: secretHash}
	return cmd.runCommand()
}

func (cmd *extractSecretCmd) runCommand() (string, error) {
	// Loop over all pushed data from all inputs, searching for one that hashes
	// to the expected hash.  By searching through all data pushes, we avoid any
	// issues that could be caused by the initiator redeeming the participant's
	// contract with some "nonstandard" or unrecognized transaction or script
	// type.
	for _, input := range cmd.redemptionTx.TxIn {
		pushes, err := txscript.PushedData(input.SignatureScript)
		if err != nil {
			return "", err
		}
		for _, push := range pushes {
			if bytes.Equal(sha256Hash(push), cmd.secretHash) {
				fmt.Printf("Secret: %x\n", push)
				return fmt.Sprintf("%x", push), nil
			}
		}
	}
	return "", errors.New("transaction does not contain the secret")
}
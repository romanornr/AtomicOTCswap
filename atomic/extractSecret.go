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

type ExtractedSecret struct {
	Secret string `json:"secret"`
}

func ExtractSecret(redemptionTransaction string, secretHash string) (extractedSecret ExtractedSecret, err error) {
	redemptionTxBytes, err := hex.DecodeString(redemptionTransaction)
	if err != nil {
		return extractedSecret, fmt.Errorf("failed to decode redemption transaction: %v", err)
	}

	var redemptionTx wire.MsgTx
	err = redemptionTx.Deserialize(bytes.NewReader(redemptionTxBytes))
	if err != nil {
		return extractedSecret, fmt.Errorf("failed to decode redemptioon transaction: %v", err)
	}

	secret, err := hex.DecodeString(secretHash)
	if err != nil {
		return extractedSecret, errors.New("secret hash must be hex encoded")
	}
	if len(secretHash) != sha256.Size {
		return extractedSecret, errors.New("secret hash wrong size")
	}

	cmd := &extractSecretCmd{redemptionTx: &redemptionTx, secretHash: secret}
	return cmd.runCommand()
}

func (cmd *extractSecretCmd) runCommand() (extractedSecret ExtractedSecret, err error) {
	// Loop over all pushed data from all inputs, searching for one that hashes
	// to the expected hash.  By searching through all data pushes, we avoid any
	// issues that could be caused by the initiator redeeming the participant's
	// contract with some "nonstandard" or unrecognized transaction or script
	// type.
	for _, input := range cmd.redemptionTx.TxIn {
		pushes, err := txscript.PushedData(input.SignatureScript)
		if err != nil {
			return extractedSecret, err
		}
		for _, push := range pushes {
			if bytes.Equal(sha256Hash(push), cmd.secretHash) {
				extractedSecret.Secret = fmt.Sprintf("%x", push)
				return extractedSecret, nil
			}
		}
	}
	return extractedSecret, errors.New("transaction does not contain the secret")
}

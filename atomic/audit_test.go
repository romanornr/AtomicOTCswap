package atomic

import (
	btcutil "github.com/viacoin/viautil"
	"testing"
)

func TestAuditContract(t *testing.T) {
	contract, err := AuditContract("6382012088a82095b1269e67e58860bb591ccec2efd34700b8fbe56b808e0e1ca66c77c3785d888876a91424cc424c1e5e977175d2b20012554d39024bd68f6704d4441e5cb17576a91424cc424c1e5e977175d2b20012554d39024bd68f6888ac", "0200000001d6af699222767cb7c62ce56bd50a070acb44024095cbfdbad780cb8c10b62023010000006b48304502210096ae3d8bc5b863e99316e6a9a94cd4194c0b2502953725262fa7c8b764e1cf16022071819d34853db65a03aaab9143f98e35f9081db92656a1dd6984f4ce5fb22683012102a7b08bb2a3609223a185761231d815e287ec13b74ccff3feb274253f7737356affffffff01a08601000000000017a914db8bf8b9c38a896e814b43e5f2e0b8f580f8d1ec8700000000").runAudit()
	if err != nil {
		t.Errorf("error decoding the contract hex and contract transaction hex")
	}

	expectedAddress := "EdAm9zSqK24KgHDZ9WmyQ9A98Ytwff5g8d"
	if contract.Address.String() != expectedAddress {
		t.Errorf("Expected contract address to be %s but instead got %s\n", expectedAddress, contract.Address.String())
	}

	expectedValue, _  := btcutil.NewAmount(0.001)
	if contract.Value != expectedValue {
		t.Errorf("Expected contract value to be %v but got %v instead\n", expectedValue, contract.Value)
	}

	expectedLockTime := int64(1545487572)
	if contract.LockTime != expectedLockTime {
		t.Errorf("Expected contract locktime to be %d but got %d instead\n", expectedLockTime, contract.LockTime)
	}

	expectedRecipientAddress := "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A"
	if expectedRecipientAddress != contract.RecipientRefundAddress.String() {
		t.Errorf("Expected contract recipient address to be %s but got %s instead\n", expectedRecipientAddress, contract.RecipientRefundAddress.String())
	}

	expectedRecipientRefundAddress := "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A"
	if contract.RecipientRefundAddress.String() != expectedRecipientRefundAddress {
		t.Errorf("Expected contract recipient refund address to be %s but got %s instead\n", expectedRecipientRefundAddress, contract.RecipientRefundAddress.String())
	}
}

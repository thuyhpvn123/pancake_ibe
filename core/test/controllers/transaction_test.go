package transaction_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	c "gitlab.com/meta-node/core/controllers"
	"gitlab.com/meta-node/core/crypto"
	pb "gitlab.com/meta-node/core/proto"
)

var secretKey []byte
var lastTransaction *pb.Transaction
var transaction *pb.Transaction
var accountState *pb.AccountState

var callSmartContractChargeTransaction *pb.Transaction
var smartContractChargeFeeAccountState *pb.AccountState

func init() {
	/* load test data */
	crypto.Init()
	secretKey, _ = hex.DecodeString("1f1a8e1eadaf03a4bd7c52c66aa139236ffd324b0e9858d6f6476ff0d17b8b08")
	hash, _ := hex.DecodeString("07d50259a5d1e679b88eac052a410140c6681723e4d2b1bbe67bf8863a7b4247")
	pubKey, _ := hex.DecodeString("b9fbb03f3e14db34b502efb044062c4614677a8d9365e408e7ef72cc3adaae333cc3106f20bfee6e5747209926ab0f38")
	address, _ := hex.DecodeString("820670e501b4cb3403514e5d4aeb4db69a6389bd")
	toAddress, _ := hex.DecodeString("d142c6c25e9d2cb439cbb40cd7dc468f3858be65")
	pendingBalance := uint256.NewInt(10)
	pendingUse := pendingBalance
	amount := uint256.NewInt(4)
	fee := uint256.NewInt(1)
	balanceAfterSend := uint256.NewInt(0).Sub(
		uint256.NewInt(0).Sub(pendingBalance, amount),
		fee,
	)

	lastTransaction = &pb.Transaction{
		Hash:                hash,
		FromAddress:         address,
		ToAddress:           toAddress,
		PubKey:              pubKey,
		PendingUse:          uint256.NewInt(0).Bytes(),
		Balance:             uint256.NewInt(0).Bytes(),
		Amount:              uint256.NewInt(0).Bytes(),
		Fee:                 uint256.NewInt(0).Bytes(),
		Data:                &pb.TransactionData{},
		PreviousTransaction: c.GetEmptyTransaction(),
	}

	transaction = &pb.Transaction{
		FromAddress:         address,
		ToAddress:           toAddress,
		PubKey:              pubKey,
		PendingUse:          pendingUse.Bytes(),
		Balance:             balanceAfterSend.Bytes(),
		Amount:              amount.Bytes(),
		Fee:                 fee.Bytes(),
		Data:                &pb.TransactionData{},
		PreviousTransaction: lastTransaction,
	}
	fmt.Printf("%v", uint256.NewInt(0).SetBytes(transaction.PendingUse))
	transaction.Hash = c.GetTransactionHash(transaction)
	transaction.Sign = crypto.Sign(secretKey, transaction.Hash)
	accountState = &pb.AccountState{
		Address:        address,
		LastHash:       lastTransaction.Hash,
		Balance:        uint256.NewInt(0).Bytes(),
		PendingBalance: pendingBalance.Bytes(),
	}

	scSecretKey, _ := hex.DecodeString("338204f7c84fcf16b0dc14d74ec661b0e41e857d37f4396eb1ceaffeb2dfa427")
	scPubkey, _ := hex.DecodeString("b63763da65d5a553e052c31c0989dee38a1da1ac4f37341234cc2a957c1475d643cc7e4f5129c06516a9e3722161ba87")
	scAddress, _ := hex.DecodeString("eaab3ac9245023cdf83b2621173ee79076c48288")
	scPendingBalance := uint256.NewInt(100000)

	commissionSign := crypto.Sign(scSecretKey, address)

	smartContractChargeFeeAccountState = &pb.AccountState{
		Address:        scAddress,
		Balance:        uint256.NewInt(0).Bytes(),
		PendingBalance: scPendingBalance.Bytes(),
		SmartContractInfo: &pb.SmartContractInfo{
			CreatorPublicKey: scPubkey,
			FeeType:          *pb.FEETYPE_SMART_CONTRACT_CHARGE_FEE.Enum(),
		},
	}

	callSmartContractChargeTransaction = &pb.Transaction{
		FromAddress: address,
		ToAddress:   toAddress,
		PubKey:      pubKey,
		PendingUse:  pendingUse.Bytes(),
		Balance:     balanceAfterSend.Bytes(),
		Amount:      amount.Bytes(),
		Fee:         fee.Bytes(),
		Data: &pb.TransactionData{
			CallData: &pb.CallData{
				CommissionSign: commissionSign,
			},
		},
		PreviousTransaction: lastTransaction,
	}
}

func TestGetEmptyTransaction(t *testing.T) {
	emptyTransaction := c.GetEmptyTransaction()
	expectedHash := "0000000000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, hex.EncodeToString(emptyTransaction.Hash), expectedHash, "The two hash should be the same.")
}

func TestVerifyTransactionFromAddress(t *testing.T) {
	verifyTransactionFromAddress := c.VerifyTransactionFromAddress(transaction)
	assert.True(t, verifyTransactionFromAddress)

	cloneTnx := *transaction
	fakeAddress, _ := hex.DecodeString("0000000000000000000000000000000000000000")
	cloneTnx.FromAddress = fakeAddress
	verifyTransactionFromAddress = c.VerifyTransactionFromAddress(&cloneTnx)
	assert.False(t, verifyTransactionFromAddress)

}

func TestGetTransactionHash(t *testing.T) {
	bHash := c.GetTransactionHash(lastTransaction)
	hexHash := hex.EncodeToString(bHash)
	expectedValue := "07d50259a5d1e679b88eac052a410140c6681723e4d2b1bbe67bf8863a7b4247"

	assert.Equal(t, hexHash, expectedValue, "The two hash should be the same.")
}

func TestVerifyTransactionHashRight(t *testing.T) {
	expectedValue := true
	verifyTransactionHashRight := c.VerifyTransactionHashRight(transaction)
	assert.Equal(t, expectedValue, verifyTransactionHashRight, "The two value should be the same.")

	cloneTnx := *transaction
	fakeHash, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	cloneTnx.Hash = fakeHash
	expectedValue = false
	verifyTransactionHashRight = c.VerifyTransactionHashRight(&cloneTnx)
	assert.Equal(t, expectedValue, verifyTransactionHashRight, "The two value should be the same.")

}

func TestVerifyLastTransactionHashRight(t *testing.T) {
	expectedValue := true
	verifyLastTransactionHashRight := c.VerifyLastTransactionHashRight(transaction)
	assert.Equal(t, expectedValue, verifyLastTransactionHashRight, "The two value should be the same.")
	cloneTnx := *transaction
	clonePrvTnx := *lastTransaction

	fakeHash, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

	cloneTnx.PreviousTransaction = &clonePrvTnx
	cloneTnx.PreviousTransaction.Hash = fakeHash
	expectedValue = false
	verifyLastTransactionHashRight = c.VerifyLastTransactionHashRight(&cloneTnx)
	assert.Equal(t, expectedValue, verifyLastTransactionHashRight, "The two value should be the same.")

}

func TestVerifyMatchLastHash(t *testing.T) {
	// match
	expectedValue := true

	verifyMatchLastHash := c.VerifyMatchLastHash(transaction, accountState)
	assert.Equal(t, expectedValue, verifyMatchLastHash, "The two value should be the same.")

	// not match
	cloneTnx := *transaction
	clonePrvTnx := *lastTransaction
	fakeHash, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	cloneTnx.PreviousTransaction = &clonePrvTnx
	cloneTnx.PreviousTransaction.Hash = fakeHash
	expectedValue = false

	verifyMatchLastHash = c.VerifyMatchLastHash(&cloneTnx, accountState)
	assert.Equal(t, expectedValue, verifyMatchLastHash, "The two value should be the same.")
}

func TestVerifyTransactionAmount(t *testing.T) {
	// match
	expectedValue := true

	verifyTransactionAmount := c.VerifyTransactionAmount(transaction)
	assert.Equal(t, expectedValue, verifyTransactionAmount, "The two value should be the same.")

	// not match
	cloneTnx := *transaction
	fakeAmount := uint256.NewInt(100000)
	cloneTnx.Amount = fakeAmount.Bytes()
	verifyTransactionAmount = c.VerifyTransactionAmount(&cloneTnx)
	expectedValue = false

	assert.Equal(t, expectedValue, verifyTransactionAmount, "The two value should be the same.")

}

func TestVerifyRightBalance(t *testing.T) {
	// match
	expectedValue := true
	verifyMatchLastBalance := c.VerifyRightBalance(
		uint256.NewInt(0).SetBytes(transaction.Balance),
		uint256.NewInt(0).SetBytes(transaction.PreviousTransaction.Balance),
		uint256.NewInt(0).SetBytes(transaction.PendingUse),
		uint256.NewInt(0).SetBytes(transaction.Fee),
	)
	assert.Equal(t, expectedValue, verifyMatchLastBalance, "The two value should be the same.")

	// not match
	cloneTnx := &*transaction
	fakeBalance := uint256.NewInt(100000)
	cloneTnx.Balance = fakeBalance.Bytes()
	verifyMatchLastBalance = c.VerifyRightBalance(
		uint256.NewInt(0).SetBytes(cloneTnx.Balance),
		uint256.NewInt(0).SetBytes(cloneTnx.PreviousTransaction.Balance),
		uint256.NewInt(0).SetBytes(cloneTnx.PendingUse),
		uint256.NewInt(0).SetBytes(cloneTnx.Fee),
	)
	expectedValue = false

	assert.Equal(t, expectedValue, verifyMatchLastBalance, "The two value should be the same.")

}

func TestVerifyMatchLastBalance(t *testing.T) {
	// match
	expectedValue := true

	verifyMatchLastBalance := c.VerifyMatchLastBalance(transaction, accountState)
	assert.Equal(t, expectedValue, verifyMatchLastBalance, "The two value should be the same.")

	// not match
	cloneTnx := *transaction
	clonePrvTnx := *lastTransaction
	fakeLastBalance := uint256.NewInt(100000)
	clonePrvTnx.Balance = fakeLastBalance.Bytes()
	cloneTnx.PreviousTransaction = &clonePrvTnx
	verifyMatchLastBalance = c.VerifyMatchLastBalance(&cloneTnx, accountState)
	expectedValue = false

	assert.Equal(t, expectedValue, verifyMatchLastBalance, "The two value should be the same.")

}

func TestVerifyTransactionPendingUse(t *testing.T) {
	// match
	expectedValue := true

	verifyPendingUse := c.VerifyTransactionPendingUse(transaction, accountState)
	assert.Equal(t, expectedValue, verifyPendingUse, "The two value should be the same.")

	// not match
	cloneTnx := *transaction
	fakePendingUse := uint256.NewInt(100000)
	cloneTnx.PendingUse = fakePendingUse.Bytes()
	verifyPendingUse = c.VerifyTransactionPendingUse(&cloneTnx, accountState)
	expectedValue = false

	assert.Equal(t, expectedValue, verifyPendingUse, "The two value should be the same.")
}

func TestVerifyCommissionSign(t *testing.T) {
	// match
	expectedValue := true

	verifyCommissionSign := c.VerifyCommissionSign(callSmartContractChargeTransaction, smartContractChargeFeeAccountState)
	assert.Equal(t, expectedValue, verifyCommissionSign, "The two value should be the same.")

	// not match
	cloneTnx := *callSmartContractChargeTransaction
	fakeCommissionSign, _ := hex.DecodeString("8880d7280ea71b84444980811ffbce416951e4924889ecad0a6c183d05f98ee3f83fc7991193e686cb522bf9779ac5880d7b00cec29772cb74e084271571dccc276ac6c835e1309a0e56477d56b0f58d90648a7713b70afd71a17e1cf3066eb2")
	cloneTnx.Data.CallData.CommissionSign = fakeCommissionSign
	verifyCommissionSign = c.VerifyCommissionSign(&cloneTnx, smartContractChargeFeeAccountState)
	expectedValue = false

	assert.Equal(t, expectedValue, verifyCommissionSign, "The two value should be the same.")
}

func TestVerifySmartContractEnoughFeeToCharge(t *testing.T) {
	// match
	expectedValue := true
	guaranteeAmount := uint256.NewInt(100)
	verifySmartContractEnoughFeeToCharge := c.VerifySmartContractEnoughFeeToCharge(callSmartContractChargeTransaction, smartContractChargeFeeAccountState, guaranteeAmount)
	assert.Equal(t, expectedValue, verifySmartContractEnoughFeeToCharge, "The two value should be the same.")

	// not match
	cloneSCAccountState := *smartContractChargeFeeAccountState
	cloneSCAccountState.PendingBalance = uint256.NewInt(10).Bytes()
	verifySmartContractEnoughFeeToCharge = c.VerifySmartContractEnoughFeeToCharge(callSmartContractChargeTransaction, &cloneSCAccountState, guaranteeAmount)
	expectedValue = false

	assert.Equal(t, expectedValue, verifySmartContractEnoughFeeToCharge, "The two value should be the same.")
}

func TestVerifyTransactionRightFee(t *testing.T) {
	feeNeeded := uint256.NewInt(1)

	// match
	expectedValue := true
	feeType := *pb.FEETYPE_SMART_CONTRACT_CHARGE_FEE.Enum()
	verifyTransactionRightFee := c.VerifyTransactionRightFee(transaction, feeNeeded, feeType)
	assert.Equal(t, expectedValue, verifyTransactionRightFee, "The two value should be the same.")

	// match
	feeType = *pb.FEETYPE_USER_CHARGE_FEE.Enum()
	verifyTransactionRightFee = c.VerifyTransactionRightFee(transaction, feeNeeded, feeType)
	assert.Equal(t, expectedValue, verifyTransactionRightFee, "The two value should be the same.")

	// not match
	cloneTnx := *transaction
	cloneTnx.Fee = uint256.NewInt(0).Bytes()
	verifyTransactionRightFee = c.VerifyTransactionRightFee(&cloneTnx, feeNeeded, feeType)
	expectedValue = false
	assert.Equal(t, expectedValue, verifyTransactionRightFee, "The two value should be the same.")
}

func TestVerifyTransactionSign(t *testing.T) {
	// match
	expectedValue := true
	verifyTransactionSign := c.VerifyTransactionSign(transaction)
	assert.Equal(t, expectedValue, verifyTransactionSign, "The two value should be the same.")

	// not match
	cloneTnx := *transaction
	fakeSign, _ := hex.DecodeString("8880d7280ea71b84444980811ffbce416951e4924889ecad0a6c183d05f98ee3f83fc7991193e686cb522bf9779ac5880d7b00cec29772cb74e084271571dccc276ac6c835e1309a0e56477d56b0f58d90648a7713b70afd71a17e1cf3066eb2")
	cloneTnx.Sign = fakeSign

	verifyTransactionSign = c.VerifyTransactionSign(&cloneTnx)
	expectedValue = false
	assert.Equal(t, expectedValue, verifyTransactionSign, "The two value should be the same.")
}

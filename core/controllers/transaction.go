package controllers

import (
	"encoding/hex"
	"reflect"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	ct "gitlab.com/meta-node/core/crypto"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

func GetEmptyTransaction() *pb.Transaction {
	hash, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	fromAddress, _ := hex.DecodeString("0000000000000000000000000000000000000000")
	toAddress, _ := hex.DecodeString("0000000000000000000000000000000000000000")
	pubKey, _ := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	pendingUse, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	balance, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	amount, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	fee, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

	return &pb.Transaction{
		Hash:        hash,
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		PubKey:      pubKey,
		PendingUse:  pendingUse,
		Balance:     balance,
		Amount:      amount,
		Fee:         fee,
	}
}

func GetTransactionHash(transaction *pb.Transaction) []byte {
	hashData := &pb.HashData{
		FromAddress:  transaction.FromAddress,
		ToAddress:    transaction.ToAddress,
		PubKey:       transaction.PubKey,
		PendingUse:   transaction.PendingUse,
		Balance:      transaction.Balance,
		Amount:       transaction.Amount,
		Fee:          transaction.Fee,
		Data:         transaction.Data,
		PreviousHash: transaction.PreviousTransaction.Hash,
	}
	b, _ := proto.Marshal(hashData)
	hash := crypto.Keccak256(b)
	return hash
}

func VerifyTransactionFromAddress(transaction *pb.Transaction) bool {
	rightAddress := ct.GetByteAddress(transaction.PubKey)
	return reflect.DeepEqual(rightAddress, transaction.FromAddress)
}

func VerifyTransactionHashRight(transaction *pb.Transaction) bool {
	return reflect.DeepEqual(
		transaction.Hash,
		GetTransactionHash(transaction),
	)
}

func VerifyLastTransactionHashRight(transaction *pb.Transaction) bool {
	return VerifyTransactionHashRight(transaction.PreviousTransaction)
}

func VerifyMatchLastHash(transaction *pb.Transaction, accountState *pb.AccountState) bool {
	return reflect.DeepEqual(transaction.PreviousTransaction.Hash, accountState.LastHash)
}

func VerifyRightBalance(
	balance *uint256.Int,
	lastBalance *uint256.Int,
	pendingUse *uint256.Int,
	fee *uint256.Int,
) bool {
	// I think when over flow may have error so need this function to check
	// if late can make sure that overflow never happend then can remove this func
	// balance <= lastBalance + pendingUse - fee
	maximumBalance := uint256.NewInt(0).Sub(
		uint256.NewInt(0).Add(lastBalance, pendingUse),
		fee,
	)
	return balance.Lt(maximumBalance) || balance.Eq(maximumBalance)
}

func VerifyTransactionAmount(transaction *pb.Transaction) bool {
	amount := uint256.NewInt(0).SetBytes(transaction.Amount)
	lastBalance := uint256.NewInt(0).SetBytes(transaction.PreviousTransaction.Balance)
	balance := uint256.NewInt(0).SetBytes(transaction.Balance)
	pendingUse := uint256.NewInt(0).SetBytes(transaction.PendingUse)
	fee := uint256.NewInt(0).SetBytes(transaction.Fee)

	rightAmount := uint256.NewInt(0).Sub(
		uint256.NewInt(0).Sub(
			uint256.NewInt(0).Add(lastBalance, pendingUse),
			balance,
		),
		fee,
	)

	if !VerifyRightBalance(
		balance,
		lastBalance,
		pendingUse,
		fee,
	) {
		return false
	}

	return amount.Eq(rightAmount)
}

func VerifyMatchLastBalance(transaction *pb.Transaction, accountState *pb.AccountState) bool {
	return uint256.NewInt(0).SetBytes(transaction.PreviousTransaction.Balance).Eq(
		uint256.NewInt(0).SetBytes(accountState.Balance),
	)
}

func VerifyTransactionPendingUse(transaction *pb.Transaction, accountState *pb.AccountState) bool {
	pendingUse := uint256.NewInt(0).SetBytes(transaction.PendingUse)
	pendingBalance := uint256.NewInt(0).SetBytes(accountState.PendingBalance)
	return pendingUse.Eq(pendingBalance) || pendingUse.Lt(pendingBalance)

}

func VerifyCommissionSign(transaction *pb.Transaction, smartContractAccountState *pb.AccountState) bool {
	return ct.VerifySign(smartContractAccountState.SmartContractInfo.CreatorPublicKey, transaction.Data.CallData.CommissionSign, transaction.FromAddress)
}

func VerifySmartContractEnoughFeeToCharge(transaction *pb.Transaction, smartContractAccountState *pb.AccountState, guaranteeAmount *uint256.Int) bool {
	totalBalance := uint256.NewInt(0).Add(
		uint256.NewInt(0).SetBytes(smartContractAccountState.Balance),
		uint256.NewInt(0).SetBytes(smartContractAccountState.PendingBalance),
	)
	minimunAmount := uint256.NewInt(0).Add(
		uint256.NewInt(0).SetBytes(transaction.Fee),
		guaranteeAmount,
	)
	return totalBalance.Gt(minimunAmount) || totalBalance.Eq(minimunAmount)
}

func VerifyTransactionRightFee(transaction *pb.Transaction, feeNeeded *uint256.Int, feeType pb.FEETYPE) bool {
	if feeType != *pb.FEETYPE_SMART_CONTRACT_CHARGE_FEE.Enum() {
		txnFee := uint256.NewInt(0).SetBytes(transaction.Fee)
		return txnFee.Eq(feeNeeded) || txnFee.Gt(feeNeeded)
	}

	return true
}

func VerifyTransactionSign(transaction *pb.Transaction) bool {
	return ct.VerifySign(transaction.PubKey, transaction.Sign, transaction.Hash)
}

package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

type AccountStateObject struct {
	accountState *pb.AccountState
	db           *AccountStateDB
}

func NewAccountStateObject(
	accountState *pb.AccountState,
	db *AccountStateDB,
) *AccountStateObject {
	return &AccountStateObject{
		accountState: accountState,
		db:           db,
	}
}

func (as *AccountStateObject) GetAccountState() *pb.AccountState {
	return as.accountState
}

func (as *AccountStateObject) SetAccountState(accountState *pb.AccountState) {
	as.db.journal.append(common.BytesToAddress(as.accountState.Address))
	as.accountState = accountState
}

func (as *AccountStateObject) AddBalance(amount *uint256.Int) {
	as.db.journal.append(common.BytesToAddress(as.accountState.Address))
	as.accountState.Balance = uint256.NewInt(0).Add(
		uint256.NewInt(0).SetBytes(as.accountState.Balance),
		amount,
	).Bytes()
}

func (as *AccountStateObject) SubBalance(amount *uint256.Int) {
	as.db.journal.append(common.BytesToAddress(as.accountState.Address))
	as.accountState.Balance = uint256.NewInt(0).Sub(
		uint256.NewInt(0).SetBytes(as.accountState.Balance),
		amount,
	).Bytes()
}

func (as *AccountStateObject) AddPendingBalance(amount *uint256.Int) {
	as.db.journal.append(common.BytesToAddress(as.accountState.Address))
	as.accountState.PendingBalance = uint256.NewInt(0).Add(
		uint256.NewInt(0).SetBytes(as.accountState.PendingBalance),
		amount,
	).Bytes()
}

func (as *AccountStateObject) SubPendingBalance(amount *uint256.Int) {
	as.db.journal.append(common.BytesToAddress(as.accountState.Address))
	as.accountState.PendingBalance = uint256.NewInt(0).Sub(
		uint256.NewInt(0).SetBytes(as.accountState.PendingBalance),
		amount,
	).Bytes()
}

func (as *AccountStateObject) SetLastHash(lastHash common.Hash) {
	as.db.journal.append(common.BytesToAddress(as.accountState.Address))
	as.accountState.LastHash = lastHash.Bytes()
}

func (as *AccountStateObject) SetSmartContractInfo(smartContractInfo *pb.SmartContractInfo) {
	as.db.journal.append(common.BytesToAddress(as.accountState.Address))
	as.accountState.SmartContractInfo = smartContractInfo
}

func (as *AccountStateObject) SetStorageRoot(storageRoot common.Hash) {
	as.db.journal.append(common.BytesToAddress(as.accountState.Address))
	as.accountState.StorageRoot = storageRoot.Bytes()
}

func (as *AccountStateObject) deepCopy(db *AccountStateDB) *AccountStateObject {
	stateObject := NewAccountStateObject(proto.Clone(as.accountState).(*pb.AccountState), db)
	return stateObject
}

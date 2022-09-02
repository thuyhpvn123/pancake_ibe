package state

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	es "github.com/ethereum/go-ethereum/core/state"

	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/holiman/uint256"
	pb "gitlab.com/meta-node/core/proto"
	cu "gitlab.com/meta-node/core/utilities"
	"google.golang.org/protobuf/proto"
)

type AccountStateDB struct {
	db           es.Database
	originalRoot common.Hash // The pre-state root, before any changes were made
	trie         es.Trie
	journal      *journal

	stateObjects        map[common.Address]*AccountStateObject
	stateObjectsPending map[common.Address]struct{} // State objects finalized but not yet written to the trie
	stateObjectsDirty   map[common.Address]struct{} // State objects finalized but not yet written to the trie

	dbErr error
}

var initLastHash, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

func NewAccountStateDB(root common.Hash, db es.Database) (*AccountStateDB, error) {
	tr, err := db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	return &AccountStateDB{
		db:                  db,
		trie:                tr,
		stateObjects:        make(map[common.Address]*AccountStateObject),
		stateObjectsPending: make(map[common.Address]struct{}),
		stateObjectsDirty:   make(map[common.Address]struct{}),
		journal:             newJournal(),
	}, nil
}

func (s *AccountStateDB) UpdateTrie(root common.Hash) {
	var err error
	s.trie, err = s.db.OpenTrie(root)
	cu.CheckFatalErr("AccountStateDB UpdateTrie ERR", err)

}

// setError remembers the first non-nil error it is called with.
func (s *AccountStateDB) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

func (s *AccountStateDB) Error() error {
	return s.dbErr
}

func (s *AccountStateDB) GetOrNewAccountState(addr common.Address) *AccountStateObject {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		stateObject = s.createObject(addr)
	}
	return stateObject
}

// Exist reports whether the given account address exists in the state.
// Notably this also returns true for suicided accounts.
func (s *AccountStateDB) Exist(addr common.Address) bool {
	return s.getStateObject(addr) != nil
}

func (s *AccountStateDB) getStateObject(addr common.Address) *AccountStateObject {
	if obj := s.stateObjects[addr]; obj != nil {
		return obj
	}

	enc, err := s.trie.TryGet(addr.Bytes())
	if err != nil {
		s.setError(fmt.Errorf("getStateObject (%x) error: %v", addr.Bytes(), err))
		return nil
	}
	if len(enc) == 0 {
		return nil
	}
	accountState := &pb.AccountState{}
	err = proto.Unmarshal(enc, accountState)
	if err != nil {
		s.setError(fmt.Errorf("getStateObject (%x) error: %v", addr.Bytes(), err))
		return nil
	}
	accountObject := NewAccountStateObject(accountState, s)
	s.setStateObject(accountObject)
	return accountObject
}

func (s *AccountStateDB) setStateObject(object *AccountStateObject) {
	s.stateObjects[common.BytesToAddress(object.accountState.Address)] = object
}

func (s *AccountStateDB) createObject(addr common.Address) *AccountStateObject {
	accountState := &pb.AccountState{
		Address:        addr.Bytes(),
		LastHash:       initLastHash,
		Balance:        uint256.NewInt(0).Bytes(),
		PendingBalance: uint256.NewInt(0).Bytes(),
	}
	accountObject := NewAccountStateObject(accountState, s)
	s.setStateObject(accountObject)
	return accountObject
}

// Carrying over the balance ensures that Ether doesn't disappear.
func (s *AccountStateDB) CreateAccount(addr common.Address) {
	s.createObject(addr)
}

func (s *AccountStateDB) SetAccountState(addr common.Address, accountState *pb.AccountState) {
	stateObject := s.GetOrNewAccountState(addr)
	if stateObject != nil {
		stateObject.SetAccountState(accountState)
	}
}

func (s *AccountStateDB) AddBalance(addr common.Address, amount *uint256.Int) {
	stateObject := s.GetOrNewAccountState(addr)
	if stateObject != nil {
		stateObject.AddBalance(amount)
	}
}

func (s *AccountStateDB) SubBalance(addr common.Address, amount *uint256.Int) {
	stateObject := s.GetOrNewAccountState(addr)
	if stateObject != nil {
		stateObject.SubBalance(amount)
	}
}

func (s *AccountStateDB) AddPendingBalance(addr common.Address, amount *uint256.Int) {
	stateObject := s.GetOrNewAccountState(addr)
	if stateObject != nil {
		stateObject.AddPendingBalance(amount)
	}
}

func (s *AccountStateDB) SubPendingBalance(addr common.Address, amount *uint256.Int) {
	stateObject := s.GetOrNewAccountState(addr)
	if stateObject != nil {
		stateObject.SubPendingBalance(amount)
	}
}

func (s *AccountStateDB) SetLastHash(addr common.Address, lastHash common.Hash) {
	stateObject := s.GetOrNewAccountState(addr)
	if stateObject != nil {
		stateObject.SetLastHash(lastHash)
	}
}

func (s *AccountStateDB) SetSmartContractInfo(addr common.Address, smartContractInfo *pb.SmartContractInfo) {
	stateObject := s.GetOrNewAccountState(addr)
	if stateObject != nil {
		stateObject.SetSmartContractInfo(smartContractInfo)
	}
}

func (s *AccountStateDB) SetStorageRoot(addr common.Address, storageRoot common.Hash) {
	stateObject := s.GetOrNewAccountState(addr)
	if stateObject != nil {
		stateObject.SetStorageRoot(storageRoot)
	}
}

// Finalise finalises the state by removing the s destructed objects and clears
// the journal as well as the refunds. Finalise, however, will not push any updates
// into the tries just yet. Only IntermediateRoot or Commit will do that.
func (s *AccountStateDB) Finalise(deleteEmptyObjects bool) {
	for addr := range s.journal.dirties {
		_, exist := s.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `s.journal.dirties` but not in `s.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}

		s.stateObjectsPending[addr] = struct{}{}
		s.stateObjectsDirty[addr] = struct{}{}

		// At this point, also ship the address off to the precacher. The precacher
		// will start loading tries, and when the change is eventually committed,
		// the commit-phase will be a lot faster
	}
	// Invalidate journal because reverting across transactions is not allowed.
	s.clearJournal()
}

// IntermediateRoot computes the current root hash of the state trie.
// It is called in between transactions to get the root hash that
// goes into transaction receipts.
func (s *AccountStateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	// Finalise all the dirty storage states and write them into the tries
	s.Finalise(deleteEmptyObjects)

	usedAddrs := make([][]byte, 0, len(s.stateObjectsPending))
	for addr := range s.stateObjectsPending {
		obj := s.stateObjects[addr]
		s.updateStateObject(obj)
		usedAddrs = append(usedAddrs, common.CopyBytes(addr[:])) // Copy needed for closure
	}

	if len(s.stateObjectsPending) > 0 {
		s.stateObjectsPending = make(map[common.Address]struct{})
	}
	return s.trie.Hash()
}

func (s *AccountStateDB) clearJournal() {
	if len(s.journal.dirties) > 0 {
		s.journal = newJournal()
	}
}

// updateStateObject writes the given object to the trie.
func (s *AccountStateDB) updateStateObject(obj *AccountStateObject) {
	addr := obj.accountState.Address
	data, _ := proto.Marshal(obj.accountState)
	if err := s.trie.TryUpdate(addr, data); err != nil {
		s.setError(fmt.Errorf("updateStateObject (%x) error: %v", addr[:], err))
	}
}

// Commit writes the state to the underlying in-memory trie database.
func (s *AccountStateDB) Commit(deleteEmptyObjects bool) (common.Hash, error) {
	if s.dbErr != nil {
		return common.Hash{}, fmt.Errorf("commit aborted due to earlier error: %v", s.dbErr)
	}
	// Finalize any pending changes and merge everything into the tries
	s.IntermediateRoot(deleteEmptyObjects)

	if len(s.stateObjectsDirty) > 0 {
		s.stateObjectsDirty = make(map[common.Address]struct{})
	}

	// The onleaf func is called _serially_, so we can reuse the same account
	// for unmarshalling every time.
	root, _, err := s.trie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent common.Hash) error {
		return nil
	})
	if err != nil {
		return common.Hash{}, err
	}

	// save to disk storage
	s.db.TrieDB().Commit(root, false, nil)
	return root, err
}

func (s *AccountStateDB) Copy() *AccountStateDB {
	// Copy all the basic fields, initialize the memory ones
	state := &AccountStateDB{
		db:                  s.db,
		trie:                s.db.CopyTrie(s.trie),
		stateObjects:        make(map[common.Address]*AccountStateObject, len(s.journal.dirties)),
		stateObjectsPending: make(map[common.Address]struct{}, len(s.stateObjectsPending)),
		stateObjectsDirty:   make(map[common.Address]struct{}, len(s.journal.dirties)),
		journal:             newJournal(),
	}
	// Copy the dirty states, logs, and preimages
	for addr := range s.journal.dirties {
		// As documented [here](https://github.com/ethereum/go-ethereum/pull/16485#issuecomment-380438527),
		// and in the Finalise-method, there is a case where an object is in the journal but not
		// in the stateObjects: OOG after touch on ripeMD prior to Byzantium. Thus, we need to check for
		// nil
		if object, exist := s.stateObjects[addr]; exist {
			// Even though the original object is dirty, we are not copying the journal,
			// so we need to make sure that anyside effect the journal would have caused
			// during a commit (or similar op) is already applied to the copy.
			state.stateObjects[addr] = object.deepCopy(state)

			state.stateObjectsDirty[addr] = struct{}{}   // Mark the copy dirty to force internal (code/state) commits
			state.stateObjectsPending[addr] = struct{}{} // Mark the copy pending to force external (account) commits
		}
	}
	// Above, we don't copy the actual journal. This means that if the copy is copied, the
	// loop above will be a no-op, since the copy's journal is empty.
	// Thus, here we iterate over stateObjects, to enable copies of copies
	for addr := range s.stateObjectsPending {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = s.stateObjects[addr].deepCopy(state)
		}
		state.stateObjectsPending[addr] = struct{}{}
	}
	for addr := range s.stateObjectsDirty {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = s.stateObjects[addr].deepCopy(state)
		}
		state.stateObjectsDirty[addr] = struct{}{}
	}

	return state
}

func (s *AccountStateDB) GetLevelDbIter() ethdb.Iterator {
	return s.db.TrieDB().DiskDB().NewIterator(nil, nil)
}

func (s *AccountStateDB) GetDiskDbBatcher() ethdb.Batch {
	return s.db.TrieDB().DiskDB().NewBatch()
}

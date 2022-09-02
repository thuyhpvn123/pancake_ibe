package state

import "github.com/ethereum/go-ethereum/common"

type journal struct {
	dirties map[common.Address]int // Dirty accounts and the number of changes
}

// newJournal create a new initialized journal.
func newJournal() *journal {
	return &journal{
		dirties: make(map[common.Address]int),
	}
}

// append inserts a new modification entry to the end of the change journal.
func (j *journal) append(address common.Address) {
	j.dirties[address]++
}

package common

import (
	"github.com/ethereum/go-ethereum/common"
	pb "gitlab.com/meta-node/core/proto"
)

type NodeCheckTransactionsVote struct {
	Hash              common.Hash
	ValidTransactions *pb.Transactions
	NodeAddress       common.Address
}

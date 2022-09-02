package transactionsDB

import (
	"log"
	"os"
	"sync"

	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

type TransactionsDB struct {
	mu                 sync.Mutex
	PendingTransaction *pb.Transaction
}

var (
	instanceTransactionsDB *TransactionsDB
)

// =============== Connections funcs ================================
func GetInstanceTransactionsDB() *TransactionsDB {
	if instanceTransactionsDB == nil {
		instanceTransactionsDB = &TransactionsDB{}
	}
	return instanceTransactionsDB
}

func (txDb *TransactionsDB) SavePendingTransaction() {
	bData, _ := proto.Marshal(txDb.PendingTransaction)
	err := os.WriteFile("./data", bData, 0644)
	if err != nil {
		log.Fatalf("Error when write data %v", err)
	}
}

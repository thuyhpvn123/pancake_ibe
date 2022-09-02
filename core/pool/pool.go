package pool

import (
	"sync"

	log "github.com/sirupsen/logrus"
	pb "gitlab.com/meta-node/core/proto"
)

type Pool struct {
	mu           sync.Mutex
	transactions []*pb.Transaction
}

func NewPool(transactions []*pb.Transaction) *Pool {
	return &Pool{
		transactions: transactions,
	}
}

func (p *Pool) AddTransaction(transaction *pb.Transaction) {
	p.mu.Lock()
	log.Infof("Adding transaction to pool")
	p.transactions = append(p.transactions, transaction)
	p.mu.Unlock()
}

func (p *Pool) AddTransactions(transactions []*pb.Transaction) {
	p.mu.Lock()
	p.transactions = append(p.transactions, transactions...)
	p.mu.Unlock()
}

func (p *Pool) TakeTransaction(numberOfTransaction int) []*pb.Transaction {
	p.mu.Lock()
	if numberOfTransaction > len(p.transactions) {
		numberOfTransaction = len(p.transactions)
	}
	rs := p.transactions[len(p.transactions)-numberOfTransaction:]
	p.transactions = p.transactions[:len(p.transactions)-numberOfTransaction]
	p.mu.Unlock()
	return rs
}

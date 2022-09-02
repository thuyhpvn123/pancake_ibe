package controllers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

func GetBlockHash(block *pb.Block) common.Hash {
	blockHashData := &pb.BlockHashData{
		Count:            block.Count,
		LastEntryHash:    block.LastEntry.Hash,
		AccountStateRoot: block.AccountStateRoot,
	}
	b, err := proto.Marshal(blockHashData)
	if err != nil {
		log.Fatal(fmt.Errorf("error when marshal block hash data"))
	}
	hash := crypto.Keccak256Hash(b)
	return hash
}

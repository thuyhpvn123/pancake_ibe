package controllers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

func GetPackHash(pack *pb.Pack) common.Hash {
	packHashData := &pb.PackHashData{}
	for _, v := range pack.Transactions {
		packHashData.TransactionHashes = append(packHashData.TransactionHashes, v.Hash)
	}
	b, err := proto.Marshal(packHashData)
	if err != nil {
		log.Fatal(fmt.Errorf("error when marshal pack hash data"))
	}
	hash := crypto.Keccak256Hash(b)
	return hash
}

func ChunkPacks(packs []*pb.Pack, chunkSize int) (chunks [][]*pb.Pack) {
	cpPack := make([]*pb.Pack, len(packs))
	copy(cpPack, packs)
	for chunkSize < int(len(cpPack)) {
		chunks = append(chunks, cpPack[0:chunkSize])
		cpPack = cpPack[chunkSize:]
	}
	chunks = append(chunks, cpPack)
	return chunks
}

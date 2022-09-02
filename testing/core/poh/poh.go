package poh

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

func GenerateHash(count *uint256.Int, lastHash common.Hash, packs []*pb.Pack) common.Hash {
	packHashes := [][]byte{}

	for _, v := range packs {
		packHashes = append(packHashes, v.Hash)
	}

	hashData := &pb.PohHashData{
		PreHash:    lastHash.Bytes(),
		Count:      count.Bytes(),
		PackHashes: packHashes,
	}

	b, _ := proto.Marshal(hashData)
	return crypto.Keccak256Hash(b)
}

func GenerateEntry(
	lastCount *uint256.Int,
	lastHash common.Hash,
	numHashes int64,
	packs []*pb.Pack,
) *pb.PohEntry {
	count := lastCount
	for i := int64(0); i < numHashes-1; i++ {
		count = count.AddUint64(count, 1)
		lastHash = GenerateHash(count, lastHash, nil)
	}

	count = count.AddUint64(count, 1)
	hash := GenerateHash(count, lastHash, packs)

	return &pb.PohEntry{
		Hash:          hash.Bytes(),
		NumHashes:     numHashes,
		LastHashCount: count.Bytes(),
		Packs: &pb.Packs{
			Packs: packs,
		},
	}
}

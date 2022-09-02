package controllers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

func GetSmartContractExecuteResultHash(rs *pb.SmartContractExecuteResult) common.Hash {
	hashData := &pb.SmartContractExecuteResultHashData{
		TransactionHash:                    rs.TransactionHash,
		Type:                               rs.Type,
		MapAddressAddBalanceChange:         rs.MapAddressAddBalanceChange,
		MapAddressSubBalanceChange:         rs.MapAddressSubBalanceChange,
		MapContractAddressDeployedCodeHash: rs.MapContractAddressDeployedCodeHash,
		MapContractAddressNewStorageRoot:   rs.MapContractAddressNewStorageRoot,
		MapContractAddressNewLogHash:       rs.MapContractAddressNewLogHash,
		ExitReason:                         rs.ExitReason,
		Exception:                          rs.Exception,
		ExMsg:                              rs.ExMsg,
		Return:                             rs.Return,
	}

	b, _ := proto.Marshal(hashData)
	return crypto.Keccak256Hash(b)
}

package network

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/meta-node/client/config"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

type Connection struct {
	Address       []byte `json:"address"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	Type          string `json:"type"`
	TCPConnection net.Conn
}

func (conn Connection) SendMessage(message *pb.Message) error {
	b, err := proto.Marshal(message)
	if err != nil {
		fmt.Printf("Error when marshal %v", err)
		return err
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(b)))
	conn.TCPConnection.Write(length)
	conn.TCPConnection.Write(b)
	return nil
}

func (conn *Connection) SendInitConnection(address string) {
	bAddress, err := hex.DecodeString(address)
	if err != nil {
		fmt.Printf("Invalid address %v", err)
	}
	protoRs, _ := proto.Marshal(&pb.InitConnection{
		Address: bAddress,
		Type:    config.AppConfig.NodeType,
	})
	message := &pb.Message{
		Header: &pb.Header{
			Type:    "request",
			From:    bAddress,
			Command: "InitConnection",
		},
		Body: protoRs,
	}

	err = conn.SendMessage(message)
	if err != nil {
		fmt.Printf("Error when send started %v", err)
	}
}

func (conn *Connection) SendSubscribeToAddress(address string) {
	address = strings.ToLower(address)
	bAddress := common.FromHex(address)
	log.Infof("Subscribing to address %v", address)
	message := &pb.Message{
		Header: &pb.Header{
			Command: "SubscribeToAddress",
		},
		Body: bAddress,
	}

	conn.SendMessage(message)

}

func (conn *Connection) SendQueryLogs(q *pb.QueryLog) {
	bData, _ := proto.Marshal(q)
	message := &pb.Message{
		Header: &pb.Header{
			Command: "QueryLogs",
		},
		Body: bData,
	}

	conn.SendMessage(message)
}

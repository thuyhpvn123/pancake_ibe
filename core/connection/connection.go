package connection

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

type Connection struct {
	Address       common.Address
	IP            string
	Port          int
	TCPConnection net.Conn
	Type          string
	mu            sync.Mutex
}

func (conn *Connection) SendMessage(message *pb.Message) error {
	if conn == nil {
		return errors.New("nil conn")
	}
	conn.mu.Lock()
	b, err := proto.Marshal(message)
	if err != nil {
		fmt.Printf("Error when marshal %v", err)
		return err
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(b)))
	conn.TCPConnection.Write(length)
	conn.TCPConnection.Write(b)
	conn.mu.Unlock()
	return nil
}

func (conn *Connection) Close() {
	conn.mu.Lock()
	conn.TCPConnection.Close()
	conn.mu.Unlock()
}

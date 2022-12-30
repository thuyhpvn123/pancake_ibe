package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/core/config"
	pb "gitlab.com/meta-node/core/proto"
	"gitlab.com/meta-node/core/utilities"
	"google.golang.org/protobuf/proto"
)

type Connection struct {
	mu            sync.Mutex
	Address       common.Address
	IP            string
	Port          int
	Type          string
	TCPConnection net.Conn
	
}

func (conn *Connection) SendMessage(message *pb.Message) error {
	if conn == nil {
		return errors.New("nil conn")
	}
	conn.mu.Lock()
	defer conn.mu.Unlock()
	b, err := proto.Marshal(message)
	if err != nil {
		fmt.Printf("Error when marshal %v", err)
		return err
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(b)))
	if conn.TCPConnection == nil {
		return errors.New("nil tcp connection")
	}
	_, err = conn.TCPConnection.Write(length)
	if err != nil {
		return err
	}
	_, err = conn.TCPConnection.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (conn *Connection) Close() {
	conn.mu.Lock()
	conn.TCPConnection.Close()
	conn.mu.Unlock()
}

func getHeaderForCommand(config config.IConfig, message string) *pb.Header {
	return &pb.Header{
		Command: message,
		Pubkey:  config.GetPubkey(),
		Version: config.GetVersion(),
	}
}

func SendMessage(config config.IConfig, c *Connection, command string, pbMessage proto.Message) error {
	body := []byte{}
	if pbMessage != nil {
		body, _ = proto.Marshal(pbMessage)
	}
	message := &pb.Message{
		Header: getHeaderForCommand(config, command),
		Body:   body,
	}
	err := c.SendMessage(message)
	utilities.CheckInfoErr("Error when SendPack", err)
	return err
}

func SendBytes(config config.IConfig, c *Connection, command string, bytes []byte) {
	message := &pb.Message{
		Header: getHeaderForCommand(config, command),
		Body:   bytes,
	}
	err := c.SendMessage(message)
	utilities.CheckInfoErr("Error when SendPack", err)
}

func SendMessageWithWg(config config.IConfig, c *Connection, command string, pbMessage proto.Message, wg *sync.WaitGroup) {
	body := []byte{}
	if pbMessage != nil {
		body, _ = proto.Marshal(pbMessage)
	}
	message := &pb.Message{
		Header: getHeaderForCommand(config, command),
		Body:   body,
	}
	err := c.SendMessage(message)
	utilities.CheckInfoErr("Error when SendPack", err)
	if wg != nil {
		wg.Done()
	}
}

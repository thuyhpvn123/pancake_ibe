package network

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"google.golang.org/protobuf/proto"

	log "github.com/sirupsen/logrus"
	"gitlab.com/meta-node/client/transactionsDB"
	pb "gitlab.com/meta-node/core/proto"
)

type MessageHandler struct {
	chData chan interface{}
}

func (handler *MessageHandler) SetChan(ch chan interface{}) {
	log.Info("Set Chan")
	handler.chData = ch
}

func (handler *MessageHandler) OnConnect(conn *Connection, address string) {
	log.Info(fmt.Sprintf("OnConnect with server %s", conn.TCPConnection.RemoteAddr()))
	conn.SendInitConnection(address)
}

func (handler *MessageHandler) OnDisconnect(conn *Connection) {
	log.Info(fmt.Printf("Disconnected with server  %s, wallet address: %v", conn.TCPConnection.RemoteAddr(), conn.Address))
}

func (h MessageHandler) HandleConnection(conn *Connection) {

	for {
		bLength := make([]byte, 8)
		_, err := io.ReadFull(conn.TCPConnection, bLength)
		if err != nil {
			switch err {
			case io.EOF:
				h.OnDisconnect(conn)
				return
			default:
				h.OnDisconnect(conn)
				log.Errorf("server error: %v\n", err)
				return
			}
		}
		messageLength := binary.LittleEndian.Uint64(bLength)

		maxMsgLength := uint64(1073741824)
		if messageLength > maxMsgLength {
			log.Errorf("Invalid messageLength want less than %v, receive %v\n", maxMsgLength, messageLength)
			conn.TCPConnection.Close()
			return
		}

		data := make([]byte, messageLength)
		byteRead, err := io.ReadFull(conn.TCPConnection, data)
		if err != nil {
			switch err {
			case io.EOF:
				h.OnDisconnect(conn)
				return
			default:
				h.OnDisconnect(conn)
				log.Errorf("server error: %v\n", err)
				return
			}
		}

		if uint64(byteRead) != messageLength {
			log.Errorf("Invalid message receive byteRead !=  messageLength %v, %v\n", byteRead, messageLength)
			conn.TCPConnection.Close()
			h.OnDisconnect(conn)
			return
		}

		message := pb.Message{}

		err = proto.Unmarshal(data[:messageLength], &message)
		if err == nil {
			h.ProcessMessage(conn, &message)
		}
	}
}

func (handler *MessageHandler) ProcessMessage(conn *Connection, message *pb.Message) {
	switch message.Header.Command {
	case "InitConnection":
		handler.handleInitConnectionMessage(conn, message)
	case "ConfirmedTransaction":
		handler.handleConfirmedTransaction(conn, message)
	case "AccountState":
		handler.handlerAccountState(message)
	case "MinerGetSmartContractStateResult":
		handler.handlerSmartContractState(conn, message)
	case "TransactionError":
		handler.handlerTransactionError(message)
	case "Receipt":
		handler.handleReceipt(conn, message)
	case "NewLogs":
		handler.handleNewLogs(message)
	case "QueryLogsResult":
		handler.handleQueryLogsResult(message)
	default:
		log.Warnf("Receive invalid message %v\n", message.Header.Command)
	}
}

func (handler *MessageHandler) handleReceipt(conn *Connection, message *pb.Message) {
	log.Info("Receive Receipt from", conn.TCPConnection.RemoteAddr())
	receipt := &pb.Receipt{}
	proto.Unmarshal(message.Body, receipt)
	handler.chData <- common.Bytes2Hex(receipt.ReturnValue)
	log.Infof("Receipt: \nTransaction hash %v\nFrom %v\nTo %v\nAmount %v\nStatus %v\nReturn %v\n",
		common.BytesToHash(receipt.TransactionHash),
		common.BytesToAddress(receipt.FromAddress),
		common.BytesToAddress(receipt.ToAddress),
		uint256.NewInt(0).SetBytes(receipt.Amount),
		receipt.Status,
		common.Bytes2Hex(receipt.ReturnValue),
	)
}

func (handler *MessageHandler) handleConfirmedTransaction(conn *Connection, message *pb.Message) {
	log.Info("Receive handleConfirmedTransaction from", conn.TCPConnection.RemoteAddr())
	transaction := &pb.Transaction{}
	// save transaction to file
	bData, _ := proto.Marshal(transaction)
	err := os.WriteFile("./data_confirmed", bData, 0644)
	if err != nil {
		log.Fatalf("Error when write data %v", err)
	}
	proto.Unmarshal(message.Body, transaction)
	log.Infoln(transaction)

	transactionsDb := transactionsDB.GetInstanceTransactionsDB()
	if common.BytesToHash(transactionsDb.PendingTransaction.Hash) == common.BytesToHash(transaction.Hash) {
		transactionsDb.SavePendingTransaction()
	} else {
		log.Warn("PendingTransaction not match")
	}
}

func (handler *MessageHandler) handleInitConnectionMessage(conn *Connection, message *pb.Message) {
	log.Info("Receive InitConnection from", conn.TCPConnection.RemoteAddr())
}

func (handler *MessageHandler) handlerAccountState(message *pb.Message) {

	if len(message.Body) == 0 {
		log.Info("Account have no data")
		return
	} else {
		accountState := &pb.AccountState{}
		proto.Unmarshal(message.Body, accountState)
		select {
		case handler.chData <- accountState:
			return
		default:
		}
		log.Infof(`
			Account data: 
			Address: %v 
			lastHash:%v 
			Balance: %v 
			Pending Balance: %v 
			SmartContractInfo: %v`,
			hex.EncodeToString(accountState.Address),
			hex.EncodeToString(accountState.LastHash),
			uint256.NewInt(0).SetBytes(accountState.Balance),
			uint256.NewInt(0).SetBytes(accountState.PendingBalance),
			accountState.SmartContractInfo,
		)
		// connect to storage to get smart contract state

		if accountState.SmartContractInfo != nil {
			message := &pb.Message{
				Header: &pb.Header{
					Type:    "request",
					Command: "MinerGetSmartContractState",
				},
				Body: accountState.Address,
			}
			conn, err := net.Dial("tcp", accountState.SmartContractInfo.StorageHost)
			if err == nil {
				tcpConn := &Connection{
					TCPConnection: conn,
				}
				tcpConn.SendMessage(message)
				go handler.HandleConnection(tcpConn)
			} else {
				fmt.Print(fmt.Errorf("err when connect to storage host %v", err))
			}
		}

	}

}

func (handler *MessageHandler) handlerSmartContractState(conn *Connection, message *pb.Message) {
	log.Info("SmartContractState: ")
	if len(message.Body) == 0 {
		log.Info("Account have no data")
		return
	} else {
		rs := &pb.SmartContractStateResult{}
		proto.Unmarshal(message.Body, rs)
		fmt.Printf("code: %v\n", hex.EncodeToString(rs.SmartContractState.Code))
		fmt.Println("storage:")
		for i, v := range rs.SmartContractState.Storage {
			fmt.Printf("%v:%v\n", i, common.Bytes2Hex(v))
		}
		conn.TCPConnection.Close()
	}
}

func (handler *MessageHandler) handlerTransactionError(message *pb.Message) {
	err := &pb.Error{}
	proto.Unmarshal(message.Body, err)
	fmt.Printf("handlerTransactionError: %v\n", err)
}

type Message struct {
	Type  string
	Value interface{}
}
type EventI struct {
	Address string `json:"address"`
	Event   string `json:"name"`
	Data    string `json:"data"`
}

func (handler *MessageHandler) handleNewLogs(message *pb.Message) {
	logs := &pb.Logs{}
	proto.Unmarshal(message.Body, logs)

	fmt.Println("========== Receive NewLogs: =========== ")
	for _, log := range logs.Logs {
		address := common.Bytes2Hex(log.Address)
		data := common.Bytes2Hex(log.Data)
		// topics := common.Bytes2Hex(log.Topics)
		fmt.Println("address", address, "data", data)
		// fmt.Println("topics ne :", topics)
		handler.chData <- EventI{address, common.Bytes2Hex(log.Topics[0]), data}
		fmt.Printf("Address: %v\nData: %v\n", address, data)
		for i, t := range log.Topics {
			fmt.Printf("Topic %v: %v\n", i, common.Bytes2Hex(t))
			// handler.chData <- EventI{address, common.Bytes2Hex(log.Topics[0]), common.Bytes2Hex(t)}
		}
	}
}

func (handler *MessageHandler) handleQueryLogsResult(message *pb.Message) {
	logs := &pb.Logs{}
	var s []EventI
	proto.Unmarshal(message.Body, logs)
	fmt.Println("========== Receive Query logs result: =========== ")
	for _, log := range logs.Logs {
		address := common.Bytes2Hex(log.Address)
		data := common.Bytes2Hex(log.Data)
		topics := make(map[int]interface{})
		event := common.Bytes2Hex(log.Topics[0])
		s = append(s, EventI{address, event, data})
		fmt.Printf("Address: %v\nData: %v\n", address, data)
		for i, t := range log.Topics {
			topic := common.Bytes2Hex(t)
			topics[i] = topic
			fmt.Printf("Topic %v: %v\n", i, topic)
		}
	}
	handler.chData <- s
}

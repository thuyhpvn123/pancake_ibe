package network

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	Address               string
	IP                    string
	Port                  int
	MessageHandler        *MessageHandler
	UnInitedConnections   []Connection
	InitedConnections     map[string]Connection
	InitedConnectionsChan chan Connection
	RemoveConnectionChan  chan Connection
}

func (server *Server) Run() {
	log.Info(fmt.Sprintf("Starting server at port %d", server.Port))
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		log.Error(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
		}

		myConn := &Connection{
			TCPConnection: conn,
		}
		server.MessageHandler.OnConnect(myConn, server.Address)
		go server.MessageHandler.HandleConnection(myConn)
	}
}

func (server *Server) ConnectToParent(connection *Connection) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", connection.IP, connection.Port))
	if err != nil {
		log.Warn("Error when connect to %v:%v, wallet adress : %v", err)
	} else {
		connection.TCPConnection = conn
		server.MessageHandler.OnConnect(connection, server.Address)
		go server.MessageHandler.HandleConnection(connection)
	}
}

func ConnectToServer(serverConnectionStr string, channel chan interface{}) *Connection {
	conn, err := net.Dial("tcp", serverConnectionStr)
	if err != nil {
		panic(err)
	}

	connection := &Connection{
		TCPConnection: conn,
	}
	messageHandler := MessageHandler{channel}
	go messageHandler.HandleConnection(connection)
	return connection
}

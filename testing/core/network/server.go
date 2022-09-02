package network

type ServerHanlder struct {
	InitHandlers []interface{}
}

func NewServerHanlder() *ServerHanlder {
	var handlers []interface{}
	return &ServerHanlder{
		InitHandlers: handlers,
	}
}

// func (server *ServerHanlder) AddHandler(handler interface{}) {
// 	server.InitHandlers = append(server.InitHandlers, handler)
// }

// func (server *ServerHanlder) ProcessInitServerHandlers() {
// 	log.Info(fmt.Sprintf("ProcessInitServerHandlers %d", len(server.InitHandlers)))
// 	if len(server.InitHandlers) > 0 {
// 		for _, handler := range server.InitHandlers {
// 			go handler.(func())()
// 		}
// 	}
// }

// type Server struct {
// 	Address               string
// 	IP                    string
// 	Port                  int
// 	MessageHandler        *MessageHandler
// 	UnInitedConnections   []Connection
// 	InitedConnections     map[string]Connection
// 	InitedConnectionsChan chan Connection
// 	RemoveConnectionChan  chan Connection
// 	ServerHandler         *ServerHanlder
// }

// func (server *Server) Run() {
// 	log.Info(fmt.Sprintf("Starting server at port %d", server.Port))
// 	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port))
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	defer listener.Close()

// 	// go server.handleInitedConnectionChan()
// 	// go server.handleRemoveConnectionChan()
// 	// go server.ServerHandler.HandleSendCheckingBlock(server)
// 	// go server.ServerHandler.HandleSendCheckedBlock(server)
// 	// go server.ServerHandler.HandleAggregateCheckedBlock(server)
// 	go server.ServerHandler.ProcessInitServerHandlers()

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Error(err)
// 		}

// 		myConn := Connection{
// 			TCPConnection: conn,
// 		}
// 		server.MessageHandler.OnConnect(&myConn)
// 		go server.MessageHandler.HandleConnection(&myConn)
// 	}
// }

// func (server *Server) HandleInitedConnectionChan() {
// 	log.Info("HandleInitedConnectionChan")
// 	for {
// 		con := <-server.InitedConnectionsChan
// 		server.InitedConnections[con.Address] = con
// 		log.Info(fmt.Sprintf("Inited Connection %v", len(server.InitedConnections)))
// 	}
// }

// // func (server *Server) handleSendCheckingBlock() {

// // 	for {
// // 		// Get all miner connection and send checking block
// // 		minersConns := server.MessageHandler.Miners
// // 		var checkingBLock *pb.CheckingBlock
// // 		if len(pool.GetInstanceCheckingPool().Transactions) != 0 && len(minersConns) != 0 {
// // 			checkingBLock = GetLedgersInstance().CreateCheckingBlock(500)
// // 			for _, conn := range minersConns {
// // 				conn.SendCheckingBlock(checkingBLock)
// // 				fmt.Println("excute SendCheckingBlock", len(checkingBLock.Transactions))
// // 			}
// // 		}
// // 		time.Sleep(100000 * time.Microsecond)
// // 	}
// // }

// // func (server *Server) handleSendCheckedBlock() {

// // 	for {
// // 		// Get parent connection and send checked block
// // 		parentConnection := server.MessageHandler.NodeParent
// // 		if parentConnection != nil && len(pool.GetInstanceCheckedBlock().Transactions) != 0 {
// // 			parentConnection.SendCheckedBlock(GetLedgersInstance().CreateCheckedBlock(500))
// // 		}
// // 		time.Sleep(100000 * time.Microsecond)
// // 	}
// // }

// func RemoveUnInitedConnection(s []Connection, i int) []Connection {
// 	s[i] = s[len(s)-1]
// 	return s[:len(s)-1]
// }

// func (server *Server) HandleRemoveConnectionChan() {
// 	log.Info("HandleRemoveConnectionChan")
// 	for {
// 		con := <-server.RemoveConnectionChan
// 		delete(server.InitedConnections, con.Address)
// 		for i, v := range server.UnInitedConnections {
// 			if v == con {
// 				server.UnInitedConnections = RemoveUnInitedConnection(server.UnInitedConnections, i)
// 			}
// 		}
// 	}
// }

// func (server *Server) ConnectToServers(connections []*Connection) {
// 	for _, v := range connections {
// 		if _, ok := server.InitedConnections[v.Address]; ok {
// 			// already connected
// 			continue
// 		}

// 		conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", v.IP, v.Port))
// 		if err != nil {
// 			log.Warn("Error when connect to %v:%v, wallet adress : %v", err)
// 		} else {
// 			v.TCPConnection = conn
// 			server.MessageHandler.OnConnect(v)
// 			go server.MessageHandler.HandleConnection(v)
// 		}
// 	}
// }

// func (server *Server) ConnectToParent(connection *Connection) {
// 	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", connection.IP, connection.Port))
// 	if err != nil {
// 		log.Warn("Error when connect to %v:%v, wallet adress : %v", err)
// 	} else {
// 		connection.TCPConnection = conn
// 		server.MessageHandler.NodeParent = connection
// 		server.MessageHandler.OnConnect(connection)
// 		go server.MessageHandler.HandleConnection(connection)
// 	}
// }

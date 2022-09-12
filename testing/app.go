package main

import (
	"encoding/hex"
	// "encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/holiman/uint256"

	// "testing1/testing/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"gitlab.com/meta-node/client/config"
	"gitlab.com/meta-node/client/network"
	"gitlab.com/meta-node/client/transactionsDB"
	cc "gitlab.com/meta-node/core/controllers"
	core_crypto "gitlab.com/meta-node/core/crypto"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
)

type Message struct {
	Type string      `json:"type"`
	Msg  interface{} `json:"message"`
}
type Message1 struct {
	Type string      `json:"type"`
	Msg  []interface{} `json:"message"`
}
// type Format1 struct {
// 	Key string `json:"key"`
// 	Value string `json:"value"`
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var clients = make(map[*websocket.Conn]bool) // connected clients
var sendQueue = make(map[*websocket.Conn]chan Message)
var broadcast = make(chan interface{}) // broadcast channel
// var sendDataC  = make(chan Message)
var RouterContract = getRouterConstant()
var FactoryContract = getFactoryConstant()
var Token0Contract = getToken0Constant()
var Token1Contract = getToken1Constant()

var WbnbContract = getWbnbConstant()
var PairContract = getPairConstant()

const (
	STORAGEHOST = "61.28.238.235:3051"
)

var accounts = [...]Account{
	// {
	// 	address: "9c396e149794945b5382c66fd996e78b6ee085c1",
	// 	private: "4b88398278565510745cecbdd711ae11d330576dbc28d6f3272c62ced00249d6",
	// },
	// {
	// 	address: "d85ae9a6ef6185aea70b1b18c3d3bfd1253ea74e",
	// 	private: "28b83ded0bbad82d2b226ed482fba6965a9f3886cda3f977426eab50eeba6a92",
	// },
	{
		address: "9c987d2b2086c1171fce5c1adb568b7a9b0b1197",
		private: "666dd436598af42067fbde62cab896ffc2541c7f88c38e9f0a6d12af95f39c1f",
	},

}

type Account struct {
	address string
	private string
}
type Format struct {
	Name   string      `json:"name"`
	Hash   string      `json:"hash"`
	Format []Parameter `json:"format"`
}
type Wbnb struct {
	Address string `json:"address"`
	Approval Format `json:"approval"`
	GetApprove Format `json:"getApprove"`
}
type Factory struct {
	Address     string `json:"address"`
	PairCreated Format `json:"pairCreated"`
}
type Router struct {
	Address      string `json:"address"`
	GetPriceList Format `json:"getPriceList"`
	AddLiquidity Format `json:"addLiquidity"`
}
type Token0 struct {
	Address  string `json:"address"`
	Transfer Format `json:"transfer"`
	Approval Format `json:"approval"`
	GetApprove Format `json:"getApprove"`
	GetBalance Format `json:"GetBalance"`

}
type Token1 struct {
	Address  string `json:"address"`
	Transfer Format `json:"transfer"`
	Approval Format `json:"approval"`
	GetApprove Format `json:"getApprove"`
	GetBalance Format `json:"GetBalance"`
}
type Pair struct {
	Address  string `json:"address"`
	Transfer Format `json:"transfer"`
	Sync     Format `json:"sync"`
	Mint     Format `json:"mint"`
	Burn     Format `json:"burn"`
	Approval Format `json:"approval"`
	GetApprove Format `json:"getApprove"`
	GetBalance Format `json:"GetBalance"`
}
type Parameter struct {
	Name string `json:"address"`
	Type string `json:"offer"`
}
var tmpl *template.Template
func main() {
	// init subscribeChain
	// welcome := "Welcome to Pancakeswap"
	connRoot := network.ConnectToServer(STORAGEHOST, broadcast)
	tmpl = template.Must(template.ParseFiles("template/index.html"))
	// defer connRoot.TCPConnection.Close()
	go subscribeChain(connRoot)
	// init listen event from subscribeChain and send to all client
	go handleSubscribeMessage()
	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/frontend/",http.StripPrefix("/frontend",fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// http.ServeFile(w, r, "./frontend/template/index.html")
		// if err := templates.ExecuteTemplate(w, "index.html", welcome); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }
			tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	http.HandleFunc("/ws", websocketHandler)

	// go HandleReceiveMessage(sendDataC)

	http.ListenAndServe(":3000", nil)

	fmt.Println("Server is running: http://localhost:3000")
	// result := <-broadcast
	// log.Println("broadcast: ", result)
}
func websocketHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true

	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Register our new client
	clients[conn] = true
	sendQueue[conn] = make(chan Message)
	go HandleReceiveMessage(conn, sendQueue[conn])
	// result1 := <-sendQueue[conn]
	// log.Println("sendQueue:",result1)
	// msgdemo := []string{"Pancakeswap Demo1"}
	conn.WriteJSON(
		
		Message{Type: "message", Msg: "Pancakeswap Demo1"})
		// Message1{Type: "message", Msg: msgdemo})

	//Get init marketplate data
	// go getEvent("", RouterContract.Address, "", "", "", "", "", "", "100", "1", conn)
	// go getEvent("", Token0Contract.Address, "", "", "", "", "", "", "100", "1", conn)
	// go getEvent("", Token1Contract.Address, "", "", "", "", "", "", "100", "1", conn)

	log.Println("Client Connected successfully") //write on server terminal
	// // Make sure we close the connection when the function returns
	// defer conn.Close()
	for {
		var msg Message

		// Read in a new message as JSON and map it to a Message object
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			delete(sendQueue, conn)
			break
		}
		// Send the newly received message to the broadcast channel
		// sendDataC <- msg

		// log message
		fmt.Println(msg)
		switch msg.Type {
		// case "GetPriceList":
		// 	//amount:="1000000000000000000"
		// 	// s := strconv.Itoa(msg.Msg)
		// 	// s := fmt.Sprintf("%d", msg.Msg.(string))
		// 	fmt.Println("This is PriceList")
		// 	// fmt.Printf(", msg:%v\n", msg.Msg.(string))
		// 	path := []string{Token0Contract.Address, Token1Contract.Address}
		// 	go getPriceList(msg.Msg.(string), path, conn)
		
		case "GetPriceList":
			fmt.Println("This is PriceList")
			fmt.Println("msg:", msg)
			msg1 := msg.Msg.(string)
			fmt.Printf("msgType:%T\n", msg1)
				amount := msg1[:strings.IndexByte(msg1, ',')]			
				pathString := msg1[strings.IndexByte(msg1, ',')+1:]
				token1:=pathString[:strings.IndexByte(pathString, ',')]
				token2:=pathString[strings.IndexByte(pathString, ',')+1:]
				path :=[]string{token1,token2}

			go getPriceList(amount, path, conn)
			//msg.Msg.(string)

		// case "AddLiquidity":
		// 	amountA := "1000000000000000000"
		// 	amountB := "1000000000000000000"
		// 	amountAmin := "1000000000000000000"
		// 	amountBmin := "1000000000000000000"
		// 	to := "9c396e149794945b5382c66fd996e78b6ee085c1"
		// 	time := "100000000000"
		// 	fmt.Println("This is Addliquidity")
		// 	go addLiquidity(Token0Contract.Address, Token1Contract.Address, amountA, amountB,
		// 		amountAmin, amountBmin, to, time, conn)
		case "GetApprove":
			fmt.Println("This is GetApprove")
			owner:= accounts[0].address
			spender :=RouterContract.Address
			go getApprove(msg.Msg.(string),owner,spender, conn)
		case "GetBalance":
			fmt.Println("This is GetBalance")
			owner:= accounts[0].address
			go getBalance(msg.Msg.(string),owner, conn)

		default:
		}
	}
}
func HandleReceiveMessage(conn *websocket.Conn, sendChan chan Message) {
	// ,
	for {
		// Grab the next message from the broadcast channel
		msg := <-sendChan

		// Send it out to every client that is currently connected
		for client := range clients {
			// Duy added

			err := client.WriteJSON(msg)
			log.Println(msg)
			// err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				// client.Close()
				conn.Close()
				delete(clients, conn)
				// delete(sendQueue,conn)
			}
		}
	}
}
// func addLiquidity(token0Addr string, token1Addr string, amountA string, amountB string, amountAmin string, amountBmin string, to string, time string, ws *websocket.Conn) {
// 	account := getAccountAvailable()
// 	conn, chCallData := initCallConnection(account.address)
// 	defer conn.TCPConnection.Close()
// 	input := formatCall(RouterContract.AddLiquidity.Hash,
// 		RouterContract.AddLiquidity.Format,
// 		[]interface{}{
// 			token0Addr, token1Addr, amountA, amountB, amountAmin, amountBmin, to, time,
// 		})
// 	fmt.Println(input)
// 	sendCallData(conn, enterAddress(RouterContract.Address), common.FromHex(input), chCallData, account.private)
// 	receiver := (<-chCallData).(string)
// 	fmt.Println("receiver Add:", receiver)
// 	giveBackAccount(account)
// 	go sentToClient(ws, "AddLiquidity", receiver)
// }

func getPriceList(num string, path []string, ws *websocket.Conn) {
	account := getAccountAvailable()
	conn, chCallData := initCallConnection(account.address)
	defer conn.TCPConnection.Close()
	input := formatCall(RouterContract.GetPriceList.Hash,
		RouterContract.GetPriceList.Format,
		[]interface{}{
			num,
			path,
		})
	sendCallData(conn, enterAddress(RouterContract.Address), common.FromHex(input), chCallData, account.private)
	receiver := (<-chCallData).(string)
	// fmt.Println("receiver GetPriceList:", receiver)
	giveBackAccount(account)
	go sentToClient(ws, "GetPriceList", formatNumber(receiver[192:]))
}
func getApprove(tokenAdr string ,owner string, spender string, ws *websocket.Conn) {
	account := getAccountAvailable()
	conn, chCallData := initCallConnection(account.address)
	defer conn.TCPConnection.Close()
	fmt.Println("tokenAdr:",tokenAdr)
	fmt.Println("PairAdd:",PairContract.Address)
	switch strings.ToLower(tokenAdr){
	case strings.ToLower(Token0Contract.Address):
		input := formatCall(Token0Contract.GetApprove.Hash,
			Token0Contract.GetApprove.Format,
			[]interface{}{
				owner,
				spender,
			})
		// fmt.Println("input:", input)
		sendCallData(conn, enterAddress(Token0Contract.Address), common.FromHex(input), chCallData, account.private)
		receiver := (<-chCallData).(string)
		// fmt.Println("receiver GetApprove:", receiver)
		giveBackAccount(account)
		go sentToClient(ws, "GetApprove", formatNumber(receiver))
	case strings.ToLower(Token1Contract.Address):
		input := formatCall(Token1Contract.GetApprove.Hash,
			Token1Contract.GetApprove.Format,
			[]interface{}{
				owner,
				spender,
			})
		// fmt.Println("input:", input)
		sendCallData(conn, enterAddress(Token1Contract.Address), common.FromHex(input), chCallData, account.private)
		receiver := (<-chCallData).(string)
		// fmt.Println("receiver GetApprove:", receiver)
		giveBackAccount(account)
		go sentToClient(ws, "GetApprove", formatNumber(receiver))
	case strings.ToLower(PairContract.Address):
		input := formatCall(PairContract.GetApprove.Hash,
			PairContract.GetApprove.Format,
			[]interface{}{
				owner,
				spender,
			})
		// fmt.Println("input:", input)
		sendCallData(conn, enterAddress(PairContract.Address), common.FromHex(input), chCallData, account.private)
		receiver := (<-chCallData).(string)
		// fmt.Println("receiver GetApprove:", receiver)
		giveBackAccount(account)
		go sentToClient(ws, "GetApprove", formatNumber(receiver))
	case strings.ToLower(WbnbContract.Address):
		input := formatCall(WbnbContract.GetApprove.Hash,
			WbnbContract.GetApprove.Format,
			[]interface{}{
				owner,
				spender,
			})
		// fmt.Println("input:", input)
		sendCallData(conn, enterAddress(WbnbContract.Address), common.FromHex(input), chCallData, account.private)
		receiver := (<-chCallData).(string)
		// fmt.Println("receiver GetApprove:", receiver)
		giveBackAccount(account)
		go sentToClient(ws, "GetApprove", formatNumber(receiver))

	}
	
}
func getBalance(tokenAdr string ,owner string, ws *websocket.Conn) {
	account := getAccountAvailable()
	conn, chCallData := initCallConnection(account.address)
	defer conn.TCPConnection.Close()
	// fmt.Println("tokenAdr:",tokenAdr)
	// fmt.Println("token1Add:",Token1Contract.Address)
	switch strings.ToLower(tokenAdr){
	case strings.ToLower(Token0Contract.Address):
		input := formatCall(Token0Contract.GetBalance.Hash,
			Token0Contract.GetBalance.Format,
			[]interface{}{
				owner,				
			})
		// fmt.Println("input:", input)
		sendCallData(conn, enterAddress(Token0Contract.Address), common.FromHex(input), chCallData, account.private)
		receiver := (<-chCallData).(string)
		// fmt.Println("receiver GetBalance:", receiver)
		giveBackAccount(account)
		go sentToClient(ws, "GetBalance", formatNumber(receiver))
	case strings.ToLower(Token1Contract.Address):
		input := formatCall(Token1Contract.GetBalance.Hash,
			Token1Contract.GetBalance.Format,
			[]interface{}{
				owner,				
			})
		// fmt.Println("input:", input)
		sendCallData(conn, enterAddress(Token1Contract.Address), common.FromHex(input), chCallData, account.private)
		receiver := (<-chCallData).(string)
		// fmt.Println("receiver GetBalance:", receiver)
		giveBackAccount(account)
		go sentToClient(ws, "GetBalance", formatNumber(receiver))
	case strings.ToLower(PairContract.Address):
		input := formatCall(PairContract.GetBalance.Hash,
			PairContract.GetBalance.Format,
			[]interface{}{
				owner,				
			})
		// fmt.Println("input:", input)
		sendCallData(conn, enterAddress(PairContract.Address), common.FromHex(input), chCallData, account.private)
		receiver := (<-chCallData).(string)
		// fmt.Println("receiver GetBalance:", receiver)
		giveBackAccount(account)
		go sentToClient(ws, "GetBalance", formatNumber(receiver))
	}
	
	
}

func sentToClient(ws *websocket.Conn, msgType string, value interface{}) {
	// fmt.Println("log ra ne:", Message{msgType, value})
	sendQueue[ws] <- Message{msgType, value}
	// fmt.Println("Toi day 22222222222222222222222")
}

func giveBackAccount(account Account) {
	queue.Lock()
	fmt.Println("give back account ")
	queue.queue[account.address] = false
	fmt.Println("end give back account ")
	queue.Unlock()
}

func formatCall(hash string, formats []Parameter, data []interface{}) string {
	s := hash
	i := 0
	for _, param := range formats {
		s += formatToHex(param.Type, data[i])
		i++
	}
	return s
} // format To Hex
func formatAddressToHex(data string) string {
	return strings.Repeat("0", 64-len(data)) + data
}
func formatNumberToHex(data int) string {
	temp := fmt.Sprintf("%X", data)
	return strings.Repeat("0", 64-len(temp)) + temp
}
func formatArray2AddressToHex(data []string) string {
	s := ""
	for _, val := range data {
		s += formatAddressToHex(val)

	}
	offset := formatNumberToHex(64)
	length := formatNumberToHex(len(data))
	result := offset + length + s
	fmt.Println(result)
	return result
}

func formatToHex(pType string, data interface{}) string {
	switch pType {
	case "address":
		return formatAddressToHex(data.(string))
	case "num":
		t, ok := data.(int)
		if !ok {
			s, err := strconv.Atoi(data.(string))
			if err != nil {
				return ""
			}
			return formatNumberToHex(s)
		}
		return formatNumberToHex(t)
	case "array":
		return formatArray2AddressToHex(data.([]string))
	default:
		return ""
	}
}

func initCallConnection(address string) (*network.Connection, chan interface{}) {
	chCallData := make(chan interface{})
	conn := network.ConnectToServer("61.28.238.31:3011", chCallData)
	conn.SendInitConnection(address)
	return conn, chCallData
}

type QueueLock struct {
	sync.Mutex
	queue map[string]bool
}

var queue = QueueLock{queue: make(map[string]bool)}

func getAccountAvailable() Account {
	for {
		// var account Account
		queue.Lock()
		fmt.Println("Find account ")
		for _, account := range accounts {
			if !queue.queue[account.address] {
				queue.queue[account.address] = true
				fmt.Println("Use account ", account.address)
				queue.Unlock()
				fmt.Println("End find account")
				return account
			}
		}
		queue.Unlock()
		fmt.Println("End find account")
		time.Sleep(time.Second / 2)
	}
}

func subscribeChain(connRoot *network.Connection) {
	fmt.Println("Storage host:", STORAGEHOST)
	connRoot.SendSubscribeToAddress(RouterContract.Address)
	connRoot.SendSubscribeToAddress(FactoryContract.Address)
	connRoot.SendSubscribeToAddress(Token0Contract.Address)
	connRoot.SendSubscribeToAddress(Token1Contract.Address)
	connRoot.SendSubscribeToAddress(PairContract.Address)
}
func getToken0Constant() Token0 {
	return Token0{
		Address: "15353116e642971722A71F9bd8Da4CFe9a10D8E9",
		Transfer: Format{
			"Transfer",
			"ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			[]Parameter{
				// {Name: "from", Type: "address"},
				// {Name: "to", Type: "address"},
				{Name: "value", Type: "num"},
			},
		},
		Approval: Format{
			"Approval",
			"8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
			[]Parameter{
				// {Name: "owner", Type: "address"},
				// {Name: "spender", Type: "address"},
				{Name: "value", Type: "num"},
			},
		},
		GetApprove: Format{
			"GetApprove", //function getAmountOunt
			"dd62ed3e",
			[]Parameter{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
			},
		},
		GetBalance: Format{
			"GetBalance", //function getAmountOunt
			"70a08231",
			[]Parameter{
				{Name: "account", Type: "address"},
			},
		},

	}
}
func getToken1Constant() Token1 {
	return Token1{
		Address: "28817d538D4d6b5522E3b61d2a64703Df7aE6562",
		Transfer: Format{
			"Transfer",
			"ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			[]Parameter{
				// {Name: "from", Type: "address"},
				// {Name: "to", Type: "address"},
				{Name: "value", Type: "num"},
			},
		},
		Approval: Format{
			"Approval",
			"8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
			[]Parameter{
				// {Name: "owner", Type: "address"},
				// {Name: "spender", Type: "address"},
				{Name: "value", Type: "num"},
			},
		},
		GetApprove: Format{
			"GetApprove", //function getAmountOunt
			"dd62ed3e",
			[]Parameter{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
			},
		},
		GetBalance: Format{
			"GetBalance", //function getAmountOunt
			"70a08231",
			[]Parameter{
				{Name: "account", Type: "address"},
			},
		},

	}
}
func getPairConstant() Pair {
	return Pair{
		Address: "6723d77cd80175ac3cdb6429c07acb4acdbd48d9",
		Transfer: Format{
			"Transfer",
			"ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			[]Parameter{
				// {Name: "from", Type: "address"},
				// {Name: "to", Type: "address"},
				{Name: "value", Type: "num"},
			},
		},
		Sync: Format{
			"Sync",
			"1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1",
			[]Parameter{
				{Name: "reserve0", Type: "num"},
				{Name: "reserve1", Type: "num"},
			},
		},
		Mint: Format{
			"Mint",
			"4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f",
			[]Parameter{
				// {Name: "sender", Type: "address"},
				{Name: "amount0", Type: "num"},
				{Name: "amount1", Type: "num"},
			},
		},
		Burn: Format{
			"Burn",
			"dccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496",
			[]Parameter{
				// {Name: "sender", Type: "address"},
				{Name: "amount0", Type: "num"},
				{Name: "amount1", Type: "num"},
			},
		},
		Approval: Format{
			"Approval",
			"8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
			[]Parameter{
				// {Name: "owner", Type: "address"},
				// {Name: "spender", Type: "address"},
				{Name: "value", Type: "num"},
			},
		},
		GetApprove: Format{
			"GetApprove", //function getAmountOunt
			"dd62ed3e",
			[]Parameter{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
			},
		},
		GetBalance: Format{
			"GetBalance", //function getAmountOunt
			"70a08231",
			[]Parameter{
				{Name: "account", Type: "address"},
			},
		},


	}
}

func getFactoryConstant() Factory {
	return Factory{
		Address: "37B34e7E96733D64Ee16C1aDa06709C2e7AA1F86",
		PairCreated: Format{
			"PairCreated",
			"0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9",
			[]Parameter{
				{Name: "token0", Type: "address"},
				{Name: "token0", Type: "address"},
				{Name: "pair", Type: "address"},
			},
		},
	}
}
func getWbnbConstant() Wbnb {
	return Wbnb{
		Address: "2dD3d44DA575795C2B88c93eFF76820d0840cc3F",
		Approval: Format{
			"Approval",
			"8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
			[]Parameter{
				// {Name: "owner", Type: "address"},
				// {Name: "spender", Type: "address"},
				{Name: "value", Type: "num"},
			},
		},
		GetApprove: Format{
			"GetApprove", //function getAmountOunt
			"dd62ed3e",
			[]Parameter{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
			},
		},

	}
}
func getRouterConstant() Router {
	return Router{
		Address: "74A1D8a484B05Ec196036C65a81C56ffB56D67fA",
		GetPriceList: Format{
			"GetPriceList", //function getAmountOunt
			"d06ca61f",
			[]Parameter{
				{Name: "amountIn", Type: "num"},
				{Name: "path", Type: "array"},
			},
		},
		AddLiquidity: Format{
			"Addliquidity", //function getAmountOunt
			"e8e33700",
			[]Parameter{
				{Name: "tokenA", Type: "address"},
				{Name: "tokenB", Type: "address"},
				{Name: "amountA", Type: "num"},
				{Name: "amountB", Type: "num"},
				{Name: "amountAmin", Type: "num"},
				{Name: "amountBmin", Type: "num"},
				{Name: "to", Type: "address"},
				{Name: "time", Type: "num"},
			},
		},
		

	}
}
func handleSubscribeMessage() {
	for {
		fmt.Println("start handleSub")
		// capture event from chain
		msg := (<-broadcast).(network.EventI)
		fmt.Println("Event ne: ", msg)
		// handle format event
		sendData := Message{}
		switch msg.Address {
		case strings.ToLower(RouterContract.Address):
			handleRouterMessage(msg, &sendData)
		case strings.ToLower(FactoryContract.Address):
			handleFactoryMessage(msg, &sendData)
		case strings.ToLower(Token0Contract.Address):
			handleToken0Message(msg, &sendData)
			fmt.Println("handleToken0Message")
		case strings.ToLower(Token1Contract.Address):
			handleToken1Message(msg, &sendData)
			fmt.Println("handleToken1Message")
		case strings.ToLower(PairContract.Address):
			handlePairMessage(msg, &sendData)
		default:
			fmt.Println("Unknow Event: ", msg)
		}

		//   Send it out to every client that is currently connected
		for client := range clients {
			sendQueue[client] <- sendData
			// sendDataC <- sendData
			// }
		}
	}
}

type Event struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func handleRouterMessage(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case RouterContract.AddLiquidity.Hash:
		*sendData = Message{
			"From:" + RouterContract.Address,
			Event{RouterContract.AddLiquidity.Name, formatEvent(
				string(msg.Data),
				RouterContract.AddLiquidity.Format,
			)},
		}
	case RouterContract.AddLiquidity.Hash:
		*sendData = Message{
			"From:" + RouterContract.Address,
			Event{RouterContract.GetPriceList.Name, formatEvent(
				string(msg.Data),
				RouterContract.GetPriceList.Format,
			)},
		}
	default:
	}
}
func handleFactoryMessage(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case FactoryContract.PairCreated.Hash:
		*sendData = Message{
			"From:" + FactoryContract.Address,
			Event{FactoryContract.PairCreated.Name, formatEvent(
				string(msg.Data),
				FactoryContract.PairCreated.Format,
			)},
		}
	default:
	}
}
func handleToken0Message(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case Token0Contract.Transfer.Hash:
		*sendData = Message{
			"From:" + Token0Contract.Address,
			Event{Token0Contract.Transfer.Name, formatEvent(
				string(msg.Data),
				Token0Contract.Transfer.Format,
			)},
		}
	case Token0Contract.Approval.Hash:
		*sendData = Message{
			"From:" + Token0Contract.Address,
			Event{Token0Contract.Approval.Name, formatEvent(
				string(msg.Data),
				Token0Contract.Approval.Format,
			)},
		}

	default:
	}
}
func handleToken1Message(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case Token1Contract.Transfer.Hash:
		*sendData = Message{
			"From:" + Token1Contract.Address,
			Event{Token1Contract.Transfer.Name, formatEvent(
				string(msg.Data),
				Token1Contract.Transfer.Format,
			)},
		}
	case Token1Contract.Approval.Hash:
		*sendData = Message{
			"From:" + Token1Contract.Address,
			Event{Token1Contract.Approval.Name, formatEvent(
				string(msg.Data),
				Token1Contract.Approval.Format,
			)},
		}

	default:
	}
}

func handlePairMessage(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case PairContract.Transfer.Hash:
		*sendData = Message{
			"From:" + PairContract.Address,
			Event{PairContract.Transfer.Name, formatEvent(
				string(msg.Data),
				PairContract.Transfer.Format,
			)},
		}
	case PairContract.Sync.Hash:
		*sendData = Message{
			"From:" + PairContract.Address,
			Event{PairContract.Sync.Name, formatEvent(
				string(msg.Data),
				PairContract.Sync.Format,
			)},
		}
	case PairContract.Mint.Hash:
		*sendData = Message{
			"From:" + PairContract.Address,
			Event{PairContract.Mint.Name, formatEvent(
				string(msg.Data),
				PairContract.Mint.Format,
			)},
		}

	default:
	}
}

func formatEvent(data string, format []Parameter) map[string]interface{} {
	formatData := make(map[string]interface{})
	for _, item := range format {
		switch item.Type {
		case "num":
			formatData[item.Name] = formatNumber(data[:64])
		case "address":
			formatData[item.Name] = formatAddress(data[:64])
		case "bool":
			formatData[item.Name] = formatBool(data[:64])
		default:
		}
		data = data[64:]
	}
	return formatData
}

func formatAddress(s string) string {
	return "0x" + s[len(s)-40:]
}
func formatBool(s string) bool {
	return s[len(s)-1] == '1'
}
func formatNumber(s string) string { //convert from hex to int dang string
	i := new(big.Int)
	i.SetString(s, 16)
	return i.String()
}

func getEvent(transactionHash, address, fromBlock,
	toBlock, topic1, topic2, topic3, topic4, slimit, spage string,
	ws *websocket.Conn) {

	chGet := make(chan interface{})
	conn := network.ConnectToServer(STORAGEHOST, chGet)
	conn.SendInitConnection(config.AppConfig.Address)
	defer conn.TCPConnection.Close()
	q := &pb.QueryLog{
		TransactionHash: common.FromHex(transactionHash),
		Address:         common.FromHex(address),
		FromBlock:       common.FromHex(fromBlock),
		ToBlock:         common.FromHex(toBlock),
		Topic1:          common.FromHex(topic1),
		Topic2:          common.FromHex(topic2),
		Topic3:          common.FromHex(topic3),
		Topic4:          common.FromHex(topic4),
	}
	limit, _ := strconv.Atoi(slimit)
	page, _ := strconv.Atoi(spage)
	q.Limit = int32(limit)
	q.Page = int32(page)
	conn.SendQueryLogs(q)
	data := (<-chGet).([]network.EventI)
	go handleGetMessage(data, ws)
}

func handleGetMessage(data []network.EventI, ws *websocket.Conn) {
	for _, item := range data {
		sendData := Message{}
		switch item.Address {
		case strings.ToLower(RouterContract.Address):
			handleRouterMessage(item, &sendData)
		case strings.ToLower(FactoryContract.Address):
			handleFactoryMessage(item, &sendData)
		case strings.ToLower(Token0Contract.Address):
			handleToken0Message(item, &sendData)
		case strings.ToLower(Token1Contract.Address):
			handleToken1Message(item, &sendData)
		}
		sendQueue[ws] <- sendData
	}
}
func sendCallData(
	conn *network.Connection,
	toAddress common.Address,
	input []byte,
	chCallData chan interface{},
	private string) {
	amount := uint256.NewInt(0)
	transferFee := uint256.NewInt(1)
	// Get info from secretkey
	bRootPrikey, bRootPubkey, bRootAddress := core_crypto.GenerateKeyPairFromSecretKey(private)
	rootAddress := common.BytesToAddress(bRootAddress)
	sendGetAccountState(conn, rootAddress)
	accountInfo := (<-chCallData).(*pb.AccountState)
	// Get last transaction
	fmt.Println("accountInfo:",accountInfo)
	lastTransaction := cc.GetEmptyTransaction()
	if hex.EncodeToString(lastTransaction.Hash) != hex.EncodeToString(accountInfo.LastHash) {
		lastTransaction = readDataLastHash(hex.EncodeToString(accountInfo.LastHash))
	}

	lastTransaction.Balance = accountInfo.Balance
	lastBalance := uint256.NewInt(0).SetBytes(lastTransaction.Balance)
	u256PendingUse := uint256.NewInt(0).SetBytes(accountInfo.PendingBalance)
	balance := uint256.NewInt(0).Add(lastBalance, u256PendingUse)
	balance = uint256.NewInt(0).Sub(balance, transferFee)
	smartContractData := &pb.TransactionData{
		Type: *pb.EXECUTE_TYPE_CALL.Enum(),
		CallData: &pb.CallData{
			Input: input,
		},
	}
	smartContractData.CallData.CommissionSign = []byte{}

	txt := &pb.Transaction{
		FromAddress:         bRootAddress,
		ToAddress:           toAddress.Bytes(),
		PubKey:              bRootPubkey,
		PendingUse:          u256PendingUse.Bytes(),
		Balance:             balance.Bytes(),
		Amount:              amount.Bytes(),
		Fee:                 transferFee.Bytes(),
		Data:                smartContractData,
		PreviousTransaction: lastTransaction,
	}
	hash := cc.GetTransactionHash(txt)
	txt.Hash = hash
	// Sign and VerifyTransaction
	txt.Sign = core_crypto.Sign(bRootPrikey, hash)
	cc.VerifyTransactionAmount(txt)

	txtB, err := proto.Marshal(txt)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(fmt.Sprintf("./datas/%s", hex.EncodeToString(txt.Hash)), txtB, 0644)
	if err != nil {
		log.Fatal(err)
	}

	message := &pb.Message{
		Header: &pb.Header{
			Type:    "request",
			From:    bRootAddress,
			Command: "SendTransaction",
		},
		Body: txtB,
	}
	conn.SendMessage(message)
	fmt.Println("sendMessage")
	txt.PreviousTransaction = &pb.Transaction{
		Hash: txt.PreviousTransaction.Hash,
	}
	transactionsDB.GetInstanceTransactionsDB().PendingTransaction = txt
	fmt.Printf("Transaction sent")
}
func readDataLastHash(hash string) *pb.Transaction {
	dat, err := os.ReadFile(fmt.Sprintf("./datas/%s", hash))
	if err != nil {
		log.Fatalf("Error when write data %v", err)
	}
	transaction := &pb.Transaction{}
	proto.Unmarshal(dat, transaction)
	return transaction
}

func sendGetAccountState(
	parentConn *network.Connection,
	address common.Address,
) {
	message := &pb.Message{
		Header: &pb.Header{
			Type:    "request",
			Command: "GetAccountState",
		},
		Body: address.Bytes(),
	}
	parentConn.SendMessage(message)
	fmt.Println("Sended Account Info")
}
func enterAddress(message string) common.Address {
	address := strings.Replace(message, "\n", "", -1)
	address = strings.ToLower(message)
	return common.HexToAddress(address)
}

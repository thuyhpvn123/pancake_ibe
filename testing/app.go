package main

import (
	"context"
	"encoding/hex"
	// "math"
	"encoding/json"
	"flag"
	"reflect"
	"gitlab.com/meta-node/client/network"
	"gitlab.com/meta-node/client/network/messages"
	"fmt"
	"html/template"
	. "github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/holiman/uint256"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.com/meta-node/client/config"
	"gitlab.com/meta-node/client/transactionsDB"
	cc "gitlab.com/meta-node/core/controllers"
	core_crypto "gitlab.com/meta-node/core/crypto"
	cn "gitlab.com/meta-node/core/network"
	pb "gitlab.com/meta-node/core/proto"
	"google.golang.org/protobuf/proto"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Type string      `json:"type"`
	Msg  interface{} `json:"message"`
}
type Message1 struct {
	Type string        `json:"type"`
	Msg  []interface{} `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
type CallData struct {
	server *Server
	client *Client
}
type Client struct {
	ws     *websocket.Conn
	server *Server
	caller CallData
	sync.Mutex
	sendChan chan Message
}
type Account struct {
	Address string
	Private string
}
type Server struct {
	sync.Mutex
	clients           ClientList
	broadcast         chan interface{}
	database          Database
	subscribe         Subscribe
	availableAccounts chan Account
	contractABI       map[string]*ContractABI
	config            config.Config
}
type Subscribe struct {
	server        *Server
	subscribeChan chan interface{}
	connRoot      *cn.Connection
	// subscribeConn *cn.Connection
}
type ContractABI struct {
	Name    string
	Address string
	Abi     ABI
}
type Database struct {
	sync.RWMutex
	data []Message
}
type ClientList struct {
	sync.RWMutex
	data map[*websocket.Conn]Client
}
type Contract struct {
	Name    string
	Address string
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var sendQueue = make(map[*websocket.Conn]chan Message)
var broadcast = make(chan interface{}) // broadcast channel
var RouterContract = getRouterConstant()
var FactoryContract = getFactoryConstant()
var WbnbContract = getWbnbConstant()
var PairContract = getPairConstant()
var defaultRelatedAddress [][]byte

const (
	STORAGEHOST = "35.196.167.172:3051"
)

var accounts = [...]Account{
	{
		Address: "aa39344b158f4004cac70bb4ace871a9b54baa1e",
		Private: "5808195a0d285c98dd942b7602f180ebbaa57ba15622786147a924d8e29daf4a",
	},

}
var contracts = [...]Contract{
	{Name: "router", Address: "12A4438d71606d49158100145A42BB791d1EACfA"},
	{Name: "factory", Address: "Fb2A01f9838F531d54a65F9D3DC9Fe3740A2136c"},
	{Name: "lp", Address: "374b073c39741e8022ea217c2e4e46e0176413e5"},
	{Name: "token0", Address: "B43DD2c18BbCCcE96d720A8a94Ee6e3F503BAE9C"},
	{Name: "token1", Address: "65035d3851AeDc3c7971f232eFaE29aAc641f8b5"},
	{Name: "wbnb", Address: "2a955DBfa8ac000584Fd5c006af60B23F0888F6e"},
	{Name: "cake", Address: "A913eFF3367c1a48E307Df1D23aAfF44AA8646e8"},
	{Name: "syrup", Address: "315Cb987DC3295125DAE36a9Ec6B076618D0b9Ba"},
	// {Name: "masterchef", Address: "57b32E949bf59DDE8E00d1B4F5f1163dc53277B0"},
	{Name: "mtcv2", Address: "0EE5f1926C29B6e02d151A31A31cabBF15001b32"},
	// {Name: "cakepool", Address: "5Cf52D9c1e4810b9862bDCBbF132AfeE6Eb4De69"},
	{Name: "dummyToken", Address: "8dA86AeCB25DcD43068b33862d79166993C04e63"},
	// {Name: "token", Address: "7DD84b6A8f2c6661a4a3f84979D0deD45072e16B"},
	{Name: "token", Address: "8E6143A15C2846471e695Fa38dD26ABB0276E172"},
	{Name: "comptroller", Address: "AaC9b1ffA601C0D723fEa2782ab203E28F3F0fEa"},
	{Name: "batDelegator", Address: "2B80f748592a3A55F882Be0c4bd4681bdE152511"},
	{Name: "daiDelegator", Address: "0Ff57c2Afc2f56878F7390AFd6a8c6b2EF4109db"},
	// {Name: "bat", Address: "B3d148293C6C03A7d192767fde0F9B474E80d423"},
	// {Name: "dai", Address: "27702fc7edF5A3b4819C020a47f4E74Ffd9070f5"},
}

type Format struct {
	Name   string      `json:"name"`
	Hash   string      `json:"hash"`
	Format []Parameter `json:"format"`
}
type Wbnb struct {
	Address    string `json:"address"`
	Approval   Format `json:"approval"`
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
type Pair struct {
	Address    string `json:"address"`
	Transfer   Format `json:"transfer"`
	Sync       Format `json:"sync"`
	Mint       Format `json:"mint"`
	Burn       Format `json:"burn"`
	Approval   Format `json:"approval"`
	GetApprove Format `json:"getApprove"`
	GetBalance Format `json:"GetBalance"`
}
type Parameter struct {
	Name string `json:"address"`
	Type string `json:"offer"`
}
type ListString struct {
	sync.Mutex
	data []string
}
type ListPool struct {
	sync.Mutex
	data []Pool
}
type ListUpdate struct {
	sync.Mutex
	data []Pool
}
type ListCompToken struct {
	sync.Mutex
	data []TokenCompound
}
type ListCompUser struct {
	sync.Mutex
	data []CompoundUser
}

type Pool struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Address    string             `json:"address,omitempty" bson:"address,omitempty"`
	APR        string           `json:"apr,omitempty" bson:"apr,omitempty"`
	Liquidity  string           `json:"liquidity,omitempty" bson:"liquidity,omitempty"`
	Multiplier string           `json:"multiplier,omitempty" bson:"multiplier,omitempty"`
	PoolID     string           `json:"poolid,omitempty" bson:"poolid,omitempty"`
}
type TokenCompound struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Address    string             `json:"address,omitempty" bson:"address,omitempty"`
	SupplyAPR  string           `json:"supplyapr,omitempty" bson:"supplyapr,omitempty"`
	BorrowAPR  string           `json:"borrowapr,omitempty" bson:"borrowapr,omitempty"`
	Liquidity  string           `json:"liquidity,omitempty" bson:"liquidity,omitempty"`
}
type CompoundUser struct {
	ID         primitive.ObjectID      `json:"_id,omitempty" bson:"_id,omitempty"`
	Address    string                  `json:"address,omitempty" bson:"address,omitempty"`
	SupplyBalance map[string]string    `json:"supplybalance,omitempty" bson:"supplybalance,omitempty"`
	BorrowBalance map[string]string     `json:"borrowbalance,omitempty" bson:"borrowbalance,omitempty"`
}

var tmpl *template.Template
var collection *mongo.Collection
var collection1 *mongo.Collection
var collection2 *mongo.Collection

func main() {
	var dir string
	flag.StringVar(&dir, "dir", ".", "testing")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://minigame:chinhtoi@cluster0.qpzfe.mongodb.net/?retryWrites=true&w=majority")
	client, _ := mongo.Connect(ctx, clientOptions)
	collection = client.Database("pancake").Collection("pancakestaking")
	collection1 = client.Database("pancake").Collection("compound")
	collection2 = client.Database("pancake").Collection("compoundUser")

	router := mux.NewRouter()
	server := Server{}
	server.Init()
	tmpl = template.Must(template.ParseFiles("template/index.html"))
	go server.subscribe.handleSubscribeMessage()
	router.PathPrefix("/frontend").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})
	router.HandleFunc("/ws", server.websocketHandler)
	router.HandleFunc("/staking", getAllPool).Methods("GET")
	router.HandleFunc("/staking/{id}", getPoolByID).Methods("GET")
	router.HandleFunc("/compoundToken/{id}", getCompoundTokenByID).Methods("GET")
	router.HandleFunc("/compoundTokenUser/{id}", getCompoundUserByID).Methods("GET")
	router.HandleFunc("/compoundToken", getAllCompToken).Methods("GET")
	router.HandleFunc("/compoundUser", getAllCompUser).Methods("GET")

	http.ListenAndServe(":3000", router)

	fmt.Println("Server is running: http://localhost:3000")
}

func getPoolByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var pool Pool
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"])
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, Pool{ID: id}).Decode(&pool)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(pool)
}

func getAllPool(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var pool []Pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var onepool Pool
		cursor.Decode(&onepool)
		pool = append(pool, onepool)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(pool)
}

func getCompoundTokenByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var compoundToken TokenCompound
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"])
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection1.FindOne(ctx, TokenCompound{ID: id}).Decode(&compoundToken)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(compoundToken)
}

func getAllCompToken(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var compoundToken []TokenCompound
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection1.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var onetoken TokenCompound
		cursor.Decode(&onetoken)
		compoundToken = append(compoundToken, onetoken)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(compoundToken)
}
func getCompoundUserByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user CompoundUser
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"])
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection2.FindOne(ctx, CompoundUser{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func getAllCompUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user []CompoundUser
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection2.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var oneuser CompoundUser
		cursor.Decode(&oneuser)
		user = append(user, oneuser)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func (server *Server) Init() {
	// init subscriber
	server.config = config.AppConfig
	server.contractABI = make(map[string]*ContractABI)
	var wg sync.WaitGroup
	for _, contract := range contracts {
		wg.Add(1)
		go server.getABI(&wg, contract)
	}
	wg.Wait()
	// connected clients
	server.clients.data = make(map[*websocket.Conn]Client)

	// broadcast channel
	server.broadcast = make(chan interface{})
	// available account map
	server.availableAccounts = make(chan Account, len(accounts))

	for _, account := range accounts {
		go server.GiveBackAccount(account)
	}

	server.subscribe = Subscribe{server: server}
	server.subscribe.initSub()
}

func (subscribe *Subscribe) initSub() {
	subscribe.subscribeChan = make(chan interface{})
	subscribe.connRoot = network.ConnectToServer(STORAGEHOST, subscribe.subscribeChan)
	subscribe.subscribeChain(subscribe.connRoot)
}

func (server *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {

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
	client := Client{ws: conn, sendChan: make(chan Message), server: server}
	client.init()
	conn.WriteJSON(

		Message{Type: "message", Msg: "Pancakeswap Demo1"})

	log.Println("Client Connected successfully") //write on server terminal
	// // Make sure we close the connection when the function returns
	// defer conn.Close()
	// client.caller = CallData{server: client.server, client: conn}
	// go client.server.database.transferToChan(client.sendChan) alo
	// log.Info("End init client")
	// add client into list to listen broadcast
	server.clients.Lock()
	server.clients.data[conn] = client
	server.clients.Unlock()
	//make sure remove client
	defer server.clients.Remove(conn)

	//listen websocket
	client.handleListen()
}
func (client *Client) init() {
	// send init message
	client.ws.WriteJSON(
		Message{Type: "message", Msg: "Welcome to Pancake"})
	client.caller = CallData{server: client.server, client: client}
	go client.handleMessage()
	go client.sendInitData()
	log.Info("End init client")
}
func (client *Client) handleMessage() {
	for {
		msg := <-client.sendChan
		log.Info(msg)
		err := client.ws.WriteJSON(msg)
		
		if err != nil {
			log.Printf("error: %v", err)
			client.ws.Close()
			client.server.clients.Remove(client.ws)
		}
	}
}
func (client *Client) sendInitData() {
	go client.server.database.transferToChan(client.sendChan)
}
func (client *Client) handleListen() {
	for {
		// Read in a new message as JSON and map it to a Message object
		var msg Message
		err := client.ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.ws.Close()
			client.server.clients.Remove(client.ws)
			break
		}
		// log message
		log.Info("Message from client: ", msg)
		client.handleCallChain(msg)
	}
}
func (clients *ClientList) Remove(ws *websocket.Conn) {
	clients.Lock()
	delete(clients.data, ws)
	clients.Unlock()
	log.Warn("Client disconnection")
}
func (client *Client) handleCallChain(msg Message) {
	// message := convertToArrayInterface(strings.Split((msg.Msg).(string), ","))
	// handle call
	switch msg.Type {

	case "GetPriceList":
		fmt.Println("This is PriceList")
		msg1 := msg.Msg.(string)
		amount := msg1[:strings.IndexByte(msg1, ',')]
		pathString := msg1[strings.IndexByte(msg1, ',')+1:]
		token1 := pathString[:strings.IndexByte(pathString, ',')]
		token2 := pathString[strings.IndexByte(pathString, ',')+1:]
		go client.caller.getPriceList(amount, token1, token2)
	case "GetApprove":
		fmt.Println("This is GetApprove")
		spender := contracts[0].Address
		msg1 := msg.Msg.(string)
		tokenAdd := msg1[:strings.IndexByte(msg1, ',')]
		addUser := msg1[strings.IndexByte(msg1, ',')+1:]
		go client.caller.getApprove(tokenAdd, addUser, spender)
	case "GetBalance":
		fmt.Println("This is GetBalance")
		msg1 := msg.Msg.(string)
		tokenAdd := msg1[:strings.IndexByte(msg1, ',')]
		addUser := msg1[strings.IndexByte(msg1, ',')+1:]
		go client.caller.getBalance(tokenAdd, addUser)

	// case "GetFarmPoolInfo":
	// 	fmt.Println("This is GetFarmPoolInfo")
	// 	go client.caller.getFarmPoolInfo(strings.ToLower(msg.Msg.(string)))
	
	// case "GetCompoundTab":
	// 	fmt.Println("This is GetCompoundTab")
	// 	go client.caller.getCompoundTab(strings.ToLower(msg.Msg.(string)))

	default:
		log.Warn("Require call not match: ", msg)
	}
}

func convertToArrayInterface(s []string) []interface{} {
	array := make([]interface{}, len(s))
	for i, v := range s {
		array[i] = v
	}
	return array
}

// function to pass all message in database to chan to send to client
func (database *Database) transferToChan(reciever chan Message) {
	database.RLock()
	for _, item := range database.data {
		reciever <- item
	}
	database.RUnlock()
}

func (server *Server) getABI(wg *sync.WaitGroup, contract Contract) {
	var temp ContractABI
	temp.initContract(contract)
	server.Lock()
	server.contractABI[contract.Name] = &temp
	server.Unlock()
	wg.Done()
}
func (contract *ContractABI) initContract(info Contract) {
	reader, err := os.Open("./abi/" + info.Name + ".json")
	if err != nil {
		log.Fatalf("Error occured while reading %s", "./abi/"+info.Name+".json")
	}
	contract.Abi, err = JSON(reader)
	if err != nil {
		log.Fatalf("Error occured while init abi %s", info.Name)
	}
	contract.Address = info.Address
	contract.Name = info.Name
	fmt.Println("Init ", info.Name)
}
func (contract *ContractABI) decode(name, data string) interface{} {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		log.Fatalf("Error occured while convert data to byte[] - Data: %s", data)
	}
	result := make(map[string]interface{})
	err = contract.Abi.UnpackIntoMap(result, name, bytes)
	if err != nil {
		log.Fatalf("Error occured while unpack %s", err)
	}
	return result
}
func (contract *ContractABI) encode(name string, args ...interface{}) []byte {
	formatedData := contract.formatPreEncode(contract.Abi.Methods[name].Inputs, args)
	data, err := contract.Abi.Pack(name, formatedData[:]...)
	if err != nil {
		log.Fatalf("Error occured while pack %s", err)
	}
	return data
}

func (contract *ContractABI) formatPreEncode(args Arguments, data []interface{}) []interface{} {
	i := 0
	temp := make([]interface{}, len(args))
	for _, arg := range args {
		temp[i] = formatData(arg.Type.String(), data[i])
		i++
	}
	return temp
}
// //tinh APY cua Compound
// func (caller *CallData) getCompoundTab(data ...interface{}) {
// //supply balance theo usd
// //borrow balance theo usd
// //max cho vay
// //balanceOfUnderlying(address)--Get the underlying balance of the `owner`
// //borrowBalanceCurrent(borrowBalanceCurrent)--Accrue interest to updated borrowIndex and then calculate account's borrow balance using the updated borrowIndex
// //totalSupply--liquidity
// 	//xóa toàn bộ database compound va compoundUser
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	filter := bson.M{}
// 	deletedResult, err := collection1.DeleteMany(ctx, filter )
// 	fmt.Println("deleted compound Token results:",deletedResult)
// 	if err != nil {
// 		fmt.Println("Delete compound Token fail")
// 		return
// 	}
// 	deletedResult1, err := collection2.DeleteMany(ctx, filter )
// 	fmt.Println("deleted compound user results:",deletedResult1)
// 	if err != nil {
// 		fmt.Println("Delete compound user fail")
// 		return
// 	}

// 	log.Info("getCompoundTab")
// 	contract := caller.server.contractABI["comptroller"]
// 	var list ListCompToken
// 	var list1 ListCompUser

// 	input := contract.encode("getAllMarkets")
// 	result := caller.tryCall(contract.Address, input)
// 	if result == "TimeOut" {
// 		log.Warn("getAllMarkets - Time Out")
// 		return
// 	}
// 	fmt.Println("result all markets", result)
// 	//lấy số lượng pool đang có
// 	allMarkets := contract.decode("getAllMarkets", result).(map[string]interface{})[""].([]common.Address)
// 	marketsLen := len(allMarkets)
// 	fmt.Printf("allMarkets type %T",allMarkets[0])
// 	supplyBalanceAr := map[string]string{}
// 	borrowBalanceAr := map[string]string{}
// 	compUser := CompoundUser{
// 		ID:         primitive.NewObjectID(),
// 		Address:    data[0].(string),
// 		SupplyBalance: supplyBalanceAr,
// 		BorrowBalance: borrowBalanceAr,
// 	}


// 	getMarketInfo := func(tokenAdd string ) {
// 		//lấy thông tin của pool
// 		supplyBalance:= caller.getBalance(tokenAdd,data[0])
// 		contract1 := caller.server.contractABI["daiDelegator"]
// 		inputT := contract1.encode("supplyRatePerBlock")
// 		result := caller.tryCall(tokenAdd, inputT)
// 		if result == "TimeOut" {
// 			log.Warn("GetPoolToken - Time Out")
// 			return
// 		}
// 		supplyRatePerBlock := contract1.decode("supplyRatePerBlock", result).(map[string]interface{})[""]
// 		fmt.Println("supplyRatePerBlock:",supplyRatePerBlock)
// 		borrowRatePerBlock := contract1.decode("borrowRatePerBlock", result).(map[string]interface{})[""]
// 		fmt.Println("borrowRatePerBlock:",borrowRatePerBlock)

// 		//lấy borrow balance của user
// 		input1 := contract1.encode("borrowBalanceStored", data[0])
// 		result1 := caller.tryCall(tokenAdd, input1)
// 		if result == "TimeOut" {
// 			log.Warn("borrowBalanceStored - Time Out")
// 			return
// 		}
// 		borrowBalanceStored := contract1.decode("borrowBalanceStored", result1).(map[string]interface{})[""]
// 		fmt.Println("borrowBalanceStored:",borrowBalanceStored)

// 		//lấy tổng borrow của 1 token
// 		input2 := contract1.encode("totalSupply")
// 		result2 := caller.tryCall(tokenAdd, input2)
// 		if result == "TimeOut" {
// 			log.Warn("totalSupply - Time Out")
// 			return
// 		}
// 		totalSupply := contract1.decode("totalSupply", result2).(map[string]interface{})[""]
// 		//lấy tỷ giá token/ctoken
// 		input3 := contract1.encode("exchangeRateStored")
// 		result3 := caller.tryCall(tokenAdd, input3)
// 		if result == "TimeOut" {
// 			log.Warn("exchangeRateStored - Time Out")
// 			return
// 		}
// 		exchangeRateStored := contract1.decode("exchangeRateStored", result3).(map[string]interface{})[""]

// 		liquidity := big.NewInt(0).Div(totalSupply.(*big.Int),exchangeRateStored.(*big.Int))
// 		fmt.Println("liquidity:",liquidity)

// 		//lay APR cua pool
// 		// Rate = vToken.supplyRatePerBlock(); // Integer
// 		// Rate = 37893566
// 		// BNB Mantissa = 1 * 10 ^ 18 (BNB has 18 decimal places)
// 		// Blocks Per Day = 20 * 60 * 24 (based on 20 blocks occurring every minute)
// 		// Days Per Year = 365
// 		// APY = (((Rate / BNB Mantissa * Blocks Per Day + 1) ^ Days Per Year - 1) * 100

// 		rateSupply:= float64(supplyRatePerBlock.(*big.Int).Int64())
// 		rateBorrow:= float64(borrowRatePerBlock.(*big.Int).Int64())

// 		mantissa := float64(math.Pow10(18))
// 		blockPerDay := float64(60*60*24)
// 		dayPerYear := float64(365)
// 		supplyApy := (math.Pow((rateSupply*blockPerDay)/mantissa+1,dayPerYear)-1)*100
// 		borrowApy := (math.Pow((rateBorrow*blockPerDay)/mantissa+1,dayPerYear)-1)*100

// 		// tạo token trong database
// 		newToken := TokenCompound{
// 			ID:         primitive.NewObjectID(),
// 			Address:    tokenAdd,
// 			SupplyAPR:  fmt.Sprintf("%f", supplyApy),
// 			BorrowAPR: fmt.Sprintf("%f", borrowApy),
// 			Liquidity:  liquidity.String(),
// 		}
		
// 		//ghi vào database
// 		ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()
// 		kq, err := collection1.InsertOne(ctx1, newToken)
// 		if err != nil {
// 			fmt.Println("create new token in data fail")
// 			return
// 		}
// 		fmt.Println("created new token in database")
// 		id := kq.InsertedID
// 		newToken.ID = id.(primitive.ObjectID)
// 		fmt.Println("created new token in database with id:",id.(primitive.ObjectID))
// 		list.Lock()
// 		list.data = append(list.data, newToken)
// 		list.Unlock()


// 		if supplyBalance.(*big.Int).Cmp(big.NewInt(0)) > 0 {
// 			supplyBalanceAr[tokenAdd]= supplyBalance.(*big.Int).String()
// 			// supplyTokens = append(supplyTokens, newToken)
// 		}
// 		if borrowBalanceStored.(*big.Int).Cmp(big.NewInt(0)) > 0  {
// 			borrowBalanceAr[tokenAdd]= borrowBalanceStored.(*big.Int).String()
// 			// borrowTokens = append(borrowTokens, newToken)
// 		}
		
		

// 	}
// 	for i := 0; int64(i) < int64(marketsLen); i++ {
// 		token:= fmt.Sprint(allMarkets[i])
// 		getMarketInfo(token)
// 	}
// 	//ghi user vào database
// 		ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()
// 		result2, err2 := collection2.InsertOne(ctx2, compUser)
// 		if err2 != nil {
// 			fmt.Println("create new user in data fail")
// 			return
// 		}
// 		fmt.Println("created new user in database")


// 		id := result2.InsertedID
// 		compUser.ID = id.(primitive.ObjectID)
// 		fmt.Println("reated new user in database with id:",id.(primitive.ObjectID))
// 		// fmt.Println("id inserted", result2)
// 		list.Lock()
// 		list1.data = append(list1.data, compUser)
// 		list.Unlock()

// 	go caller.sentToClient("GetFarmPoolInfo", "GetFarmPoolInfo")
// }
// //Update thông tin trên Tab Compound
// func (caller *CallData) getUpdateCompoundTab(data ...interface{}) {
// 		//xóa toàn bộ database compound va compoundUser
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()
	
// 		filter := bson.M{}
// 		deletedResult, err := collection1.DeleteMany(ctx, filter )
// 		fmt.Println("deleted compound Token results:",deletedResult)
// 		if err != nil {
// 			fmt.Println("Delete compound Token fail")
// 			return
// 		}
// 		deletedResult1, err := collection2.DeleteMany(ctx, filter )
// 		fmt.Println("deleted compound user results:",deletedResult1)
// 		if err != nil {
// 			fmt.Println("Delete compound user fail")
// 			return
// 		}
	
// 		log.Info("getCompoundTab")
// 		contract := caller.server.contractABI["comptroller"]
// 		var list ListCompToken
// 		var list1 ListCompUser
	
// 		input := contract.encode("getAllMarkets")
// 		result := caller.tryCall(contract.Address, input)
// 		if result == "TimeOut" {
// 			log.Warn("getAllMarkets - Time Out")
// 			return
// 		}
// 		fmt.Println("result all markets", result)
// 		//lấy số lượng pool đang có
// 		allMarkets := contract.decode("getAllMarkets", result).(map[string]interface{})[""].([]common.Address)
// 		marketsLen := len(allMarkets)
// 		fmt.Printf("allMarkets type %T",allMarkets[0])
// 		supplyBalanceAr := map[string]string{}
// 		borrowBalanceAr := map[string]string{}
// 		compUser := CompoundUser{
// 			ID:         primitive.NewObjectID(),
// 			Address:    data[0].(string),
// 			SupplyBalance: supplyBalanceAr,
// 			BorrowBalance: borrowBalanceAr,
// 		}
	
	
// 		getMarketInfo := func(tokenAdd string ) {
// 			//lấy thông tin của pool
// 			supplyBalance:= caller.getBalance(tokenAdd,data[0])
// 			contract1 := caller.server.contractABI["daiDelegator"]
// 			inputT := contract1.encode("supplyRatePerBlock")
// 			result := caller.tryCall(tokenAdd, inputT)
// 			if result == "TimeOut" {
// 				log.Warn("GetPoolToken - Time Out")
// 				return
// 			}
// 			supplyRatePerBlock := contract1.decode("supplyRatePerBlock", result).(map[string]interface{})[""]
// 			fmt.Println("supplyRatePerBlock:",supplyRatePerBlock)
// 			borrowRatePerBlock := contract1.decode("borrowRatePerBlock", result).(map[string]interface{})[""]
// 			fmt.Println("borrowRatePerBlock:",borrowRatePerBlock)
	
// 			//lấy borrow balance của user
// 			input1 := contract1.encode("borrowBalanceStored", data[0])
// 			result1 := caller.tryCall(tokenAdd, input1)
// 			if result == "TimeOut" {
// 				log.Warn("borrowBalanceStored - Time Out")
// 				return
// 			}
// 			borrowBalanceStored := contract1.decode("borrowBalanceStored", result1).(map[string]interface{})[""]
// 			fmt.Println("borrowBalanceStored:",borrowBalanceStored)
	
// 			//lấy tổng borrow của 1 token
// 			input2 := contract1.encode("totalSupply")
// 			result2 := caller.tryCall(tokenAdd, input2)
// 			if result == "TimeOut" {
// 				log.Warn("totalSupply - Time Out")
// 				return
// 			}
// 			totalSupply := contract1.decode("totalSupply", result2).(map[string]interface{})[""]
// 			//lấy tỷ giá token/ctoken
// 			input3 := contract1.encode("exchangeRateStored")
// 			result3 := caller.tryCall(tokenAdd, input3)
// 			if result == "TimeOut" {
// 				log.Warn("exchangeRateStored - Time Out")
// 				return
// 			}
// 			exchangeRateStored := contract1.decode("exchangeRateStored", result3).(map[string]interface{})[""]
	
// 			liquidity := big.NewInt(0).Div(totalSupply.(*big.Int),exchangeRateStored.(*big.Int))
// 			fmt.Println("liquidity:",liquidity)
	
// 			//lay APR cua pool
// 			// Rate = vToken.supplyRatePerBlock(); // Integer
// 			// Rate = 37893566
// 			// BNB Mantissa = 1 * 10 ^ 18 (BNB has 18 decimal places)
// 			// Blocks Per Day = 20 * 60 * 24 (based on 20 blocks occurring every minute)
// 			// Days Per Year = 365
// 			// APY = (((Rate / BNB Mantissa * Blocks Per Day + 1) ^ Days Per Year - 1) * 100
	
// 			rateSupply:= float64(supplyRatePerBlock.(*big.Int).Int64())
// 			rateBorrow:= float64(borrowRatePerBlock.(*big.Int).Int64())
	
// 			mantissa := float64(math.Pow10(18))
// 			blockPerDay := float64(60*60*24)
// 			dayPerYear := float64(365)
// 			supplyApy := (math.Pow((rateSupply*blockPerDay)/mantissa+1,dayPerYear)-1)*100
// 			borrowApy := (math.Pow((rateBorrow*blockPerDay)/mantissa+1,dayPerYear)-1)*100
	
// 			// tạo token trong database
// 			newToken := TokenCompound{
// 				ID:         primitive.NewObjectID(),
// 				Address:    tokenAdd,
// 				SupplyAPR:  fmt.Sprintf("%f", supplyApy),
// 				BorrowAPR: fmt.Sprintf("%f", borrowApy),
// 				Liquidity:  liquidity.String(),
// 			}
			
// 			//ghi vào database
// 			ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 			defer cancel()
// 			kq, err := collection1.InsertOne(ctx1, newToken)
// 			if err != nil {
// 				fmt.Println("create new token in data fail")
// 				return
// 			}
// 			fmt.Println("created new token in database")
// 			id := kq.InsertedID
// 			newToken.ID = id.(primitive.ObjectID)
// 			fmt.Println("created new token in database with id:",id.(primitive.ObjectID))
// 			list.Lock()
// 			list.data = append(list.data, newToken)
// 			list.Unlock()
	
	
// 			if supplyBalance.(*big.Int).Cmp(big.NewInt(0)) > 0 {
// 				supplyBalanceAr[tokenAdd]= supplyBalance.(*big.Int).String()
// 			}
// 			if borrowBalanceStored.(*big.Int).Cmp(big.NewInt(0)) > 0  {
// 				borrowBalanceAr[tokenAdd]= borrowBalanceStored.(*big.Int).String()
// 			}
			
			
	
// 		}
// 		for i := 0; int64(i) < int64(marketsLen); i++ {
// 			token:= fmt.Sprint(allMarkets[i])
// 			getMarketInfo(token)
// 		}
// 		//ghi user vào database
// 			ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 			defer cancel()
// 			result2, err2 := collection2.InsertOne(ctx2, compUser)
// 			if err2 != nil {
// 				fmt.Println("create new user in data fail")
// 				return
// 			}
// 			fmt.Println("created new user in database")
	
	
// 			id := result2.InsertedID
// 			compUser.ID = id.(primitive.ObjectID)
// 			fmt.Println("reated new user in database with id:",id.(primitive.ObjectID))
// 			list.Lock()
// 			list1.data = append(list1.data, compUser)
// 			list.Unlock()
	
// 		go caller.sentToClient("GetFarmPoolInfo", "GetFarmPoolInfo")
// }
	
// // Lấy thông tin address, liquidity, multiplier, ARP của các pool staking
// func (caller *CallData) getFarmPoolInfo(data ...interface{}) {
	
// 	//xóa toàn bộ database
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	filter := bson.M{}
// 	deletedResult, err := collection.DeleteMany(ctx, filter )
// 	fmt.Println("deleted results:",deletedResult)
// 	if err != nil {
// 		fmt.Println("Delete fail")
// 		return
// 	}
	
// 	contract := caller.server.contractABI["mtcv2"]
// 	log.Info("GetFarmPoolInfo")
// 	var list ListPool
// 	input := contract.encode("poolLength")
// 	result := caller.tryCall(contract.Address, input)
// 	if result == "TimeOut" {
// 		log.Warn("poolLength - Time Out")
// 		return
// 	}
// 	fmt.Println("result poollen", result)
// 	//lấy số lượng pool đang có
// 	poolLen := contract.decode("poolLength", result).(map[string]interface{})["pools"]
// 	poolLength := poolLen.(*big.Int).Int64()

// 	getPoolToken := func(i int) {
// 		//lấy thông tin của pool
// 		input := contract.encode("poolInfo", fmt.Sprintf("%d", i))
// 		receiver := caller.tryCall(contract.Address, input)
// 		if receiver == "TimeOut" {
// 			log.Warn("GetPoolInfo - Time Out")
// 			return
// 		}
// 		//lấy info multiplier của mỗi pool
// 		log.Info("GetFarmPoolInfo")
// 		allocPoint := contract.decode("poolInfo", receiver).(map[string]interface{})["allocPoint"]
// 		multiplier := big.NewInt(0).Div(allocPoint.(*big.Int), big.NewInt(100))
// 		inputT := contract.encode("lpToken", fmt.Sprintf("%d", i))
// 		receiver1 := caller.tryCall(contract.Address, inputT)
// 		if receiver == "TimeOut" {
// 			log.Warn("GetPoolToken - Time Out")
// 			return
// 		}
// 		//lấy info địa chỉ của token pool
// 		token := contract.decode("lpToken", receiver1).(map[string]interface{})[""]
// 		contract1 := caller.server.contractABI["lp"]
// 		input1 := contract1.encode("balanceOf", fmt.Sprintf("%v", contract.Address))
// 		result := caller.tryCall(fmt.Sprintf("%v", token), input1)
// 		if result == "TimeOut" {
// 			log.Warn("GetBalancePool - Time Out")
// 			return
// 		}
// 		liquidity := contract1.decode("balanceOf", result).(map[string]interface{})[""]
// 		log.Info("GetLiquidityPool - Result - ", liquidity)
// 		//lấy info liquidity của pool
// 		// list.Lock()
// 		// list.data = append(list.data, fmt.Sprint(liquidity))
// 		// list.Unlock()
// 		//lay APR cua pool
// 		var APR *big.Int
// 		cakePerblock := 1
// 		priceOfCake := big.NewInt(int64(20)) //getAmountOut in router????
// 		cakePerYear := big.NewInt(int64(cakePerblock * 60 * 60 * 24 * 365))
// 		volume24h := 50000 // dang fix cung theo usd
// 		amountOfCakePerYear := big.NewInt(0).Mul(multiplier, cakePerYear)
// 		totalValueOfCake := big.NewInt(0).Mul(amountOfCakePerYear, priceOfCake)
// 		if liquidity.(*big.Int).Cmp(big.NewInt(0)) == 0 {

// 			APR = big.NewInt(0)
// 			fmt.Println("Liquidity is zero")
// 		} else {
// 			farmBaseReward := big.NewInt(0).Div(big.NewInt(0).Mul(totalValueOfCake, new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)), (liquidity.(*big.Int))) //Percentage
// 			//farmBaseReward= totalValueOfCake)/liquidity*100/10^18
// 			totalFee := volume24h * 17 / 10000 * 365
// 			lpReward := big.NewInt(0).Div(big.NewInt(int64(totalFee)), big.NewInt(0).Div((liquidity.(*big.Int)), new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)))
// 			// lpReward = totalFee/(liquidity/10^18)*100
// 			APR = big.NewInt(0).Add(farmBaseReward, lpReward)
// 		}
// 		log.Info("GetAprPool - Result - ", APR)
// 		newPool := Pool{
// 			ID:         primitive.NewObjectID(),
// 			Address:    token.(common.Address).String(),
// 			APR:        APR.String(),
// 			Multiplier: multiplier.String(),
// 			Liquidity:  liquidity.(*big.Int).String(),
// 			PoolID : strconv.Itoa(i),
// 		}
// 		//ghi vào database
// 		ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()
// 		result1, err := collection.InsertOne(ctx1, newPool)
// 		if err != nil {
// 			fmt.Println("create newPool in data fail")
// 			return
// 		}
// 		id := result1.InsertedID
// 		newPool.ID = id.(primitive.ObjectID)
// 		fmt.Println("id.(primitive.ObjectID):",id.(primitive.ObjectID))
// 		fmt.Println("id inserted", result1)
// 		list.Lock()
// 		list.data = append(list.data, newPool)
// 		list.Unlock()

// 	}
// 	for i := 0; int64(i) < poolLength; i++ {
		
// 		getPoolToken(i)
// 	}

// 	go caller.sentToClient("GetFarmPoolInfo", list.data)
// }

//quy đổi số lượng giữa 2 token trước khi swap
func (caller *CallData) getPriceList(data ...interface{}) {

	contract := caller.server.contractABI["router"]
	log.Info("GetPriceList")
	tk0 := common.HexToAddress(data[1].(string))
	tk1 := common.HexToAddress(data[2].(string))
	path := []common.Address{tk0, tk1}
	//lấy relatedAddress là LP token
	contract1 := caller.server.contractABI["factory"]
	input1 := contract1.encode("getPair", data[1],data[2])
	result1 := caller.tryCall(contract1.Address, input1,enterRelatedAddress(""))
	if result1 == "TimeOut" {
		log.Warn("getPair - Time Out")
		return
	}
	relatedAddress := contract1.decode("getPair", result1).(map[string]interface{})[""]
	fmt.Println("relatedAddress:",relatedAddress)
	input := contract.encode("getAmountsOut", data[0], path)
	result := caller.tryCall(contract.Address, input, enterRelatedAddress(fmt.Sprint(relatedAddress)))
	if result == "TimeOut" {
		log.Warn("GetPriceList - Time Out")
		return
	}
	var out []interface{}
	price := contract.decode("getAmountsOut", result).(map[string]interface{})["amounts"]
	rv := reflect.ValueOf(price)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			out = append(out, rv.Index(i).Interface())
		}
	}
	log.Info("getPriceList - Result - ", out[1])
	go caller.sentToClient("GetPriceList", fmt.Sprint(out[1]))
}

func (caller *CallData) getApprove(data ...interface{}) {
	if strings.EqualFold(contracts[5].Address, data[0].(string)) {
		contract := caller.server.contractABI["wbnb"]
		caller.getApprove1(contract, data[1], data[2])
	} else {
		contract := caller.server.contractABI["token0"]
		caller.getApprove1(contract, data[1], data[2])
	}
}

func (caller *CallData) getApprove1(contract *ContractABI, data ...interface{}) {
	log.Info("GetApprove")
	input := contract.encode("allowance", data[0], data[1])
	result := caller.tryCall(contract.Address, input, enterRelatedAddress(""))
	if result == "TimeOut" {
		log.Warn("GetApprove - Time Out")
		return
	}
	allowance := contract.decode("allowance", result).(map[string]interface{})[""]
	log.Info("GetApprove - Result - ", allowance)
	go caller.sentToClient("GetApprove", fmt.Sprint(allowance))

}
func (caller *CallData) getBalance(data ...interface{})interface{} {
	var balance *big.Int
	if strings.EqualFold(contracts[5].Address, data[0].(string)) {
		// contract := caller.server.contractABI["wbnb"]
		balance1:=caller.getBalanceMTD(data[1])
		balance =balance1.(*uint256.Int).ToBig()
	} else {
		contract := caller.server.contractABI["token0"]
		balance=caller.getBalance1(contract, data[0], data[1]).(*big.Int)
	}
	return balance
}
func (caller *CallData) getBalance1(contract *ContractABI, data ...interface{})interface{} {
	log.Info("GetBalance")
	fmt.Println("data ne:", data)
	var list ListString
	input := contract.encode("balanceOf", data[1])
	fmt.Println(input)
	result := caller.tryCall(data[0].(string), input,enterRelatedAddress(""))
	fmt.Println(result)
	if result == "TimeOut" {
		log.Warn("GetBalance - Time Out")
		return big.NewInt(0)
	}
	price := contract.decode("balanceOf", result).(map[string]interface{})[""]
	log.Info("GetBalance - Result - ", price)
	list.Lock()
	list.data = append(list.data, fmt.Sprintf("%v", price))
	list.data = append(list.data, fmt.Sprintf("%v", data[0].(string)))
	list.Unlock()
	go caller.sentToClient("GetBalance", list.data)
	return price.(*big.Int)
}
//lấy số dư native token
func (caller *CallData) getBalanceMTD(data ...interface{})interface{} {
	log.Info("GetBalance MTD")
	account := <-caller.server.availableAccounts
	conn, chCallData := caller.initCallConnection(account.Address)
	defer conn.GetTcpConnection().Close()
	sendGetAccountState(caller.server.config, conn, common.HexToAddress(data[0].(string)))

	accountInfo := (<-chCallData).(*pb.AccountState)
	lastTransaction := cc.GetEmptyTransaction()
	lastTransaction.Balance = accountInfo.Balance
	balance := uint256.NewInt(0).SetBytes(lastTransaction.Balance)
	fmt.Println("account info:", accountInfo)
	fmt.Println("balance : %v\n", balance)
	log.Info("GetBalance MTD- Result - ", balance)
	go caller.server.GiveBackAccount(account)
	go caller.sentToClient("GetBalance", fmt.Sprint(balance))
	return balance
}

func (caller *CallData) sentToClient(msgType string, value interface{}) {
	caller.client.sendChan <- Message{msgType, value}
	sendQueue[caller.client.ws] <- Message{msgType, value}
}
func  sentToClient1(msgType string, value interface{}) {
	// sendDataC <- Message{msgType, value}
}



func (server *Server) GiveBackAccount(account Account) {
	queue.Lock()
	fmt.Println("give back account ")
	queue.queue[account.Address] = false
	fmt.Println("end give back account ")
	queue.Unlock()
	server.availableAccounts <- account
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

func (caller *CallData) initCallConnection(address string) (*cn.Connection, chan interface{}) {
	chCallData := make(chan interface{})
	conn := network.ConnectToServer("34.138.137.194:3011", chCallData)
	// conn.SendInitConnection(address)
	cn.SendMessage(caller.server.config, conn, messages.InitConnection, &pb.InitConnection{
		Address: common.FromHex(address),
		Type:    "Client",
	})

	return conn, chCallData
}


func enterRelatedAddress(address string) [][]byte {
	
	var relatedAddress [][]byte
	if address == "" {
		return defaultRelatedAddress
	}
	temp := strings.Split(address, ",")
	for _, addr := range temp {
		addressHex := common.HexToAddress(addr)
		fmt.Println("temp", addressHex)
		relatedAddress = append(relatedAddress, addressHex.Bytes())
	}
	fmt.Println("relatedAddress:",relatedAddress)
	defaultRelatedAddress = append(defaultRelatedAddress, relatedAddress...)
	return relatedAddress
}

func (caller *CallData) tryCall(address string, input []byte,relatedAddress [][]byte) string {
	i := 0
	result := "TimeOut"
	for {
		if i >= 3 {
			break
		}
		if i != 0 {
			time.Sleep(time.Second)
		}
		result = caller.call(address, input,relatedAddress)
		if result != "TimeOut" {
			log.Info("Success time - ", i)
			log.Info(" - Result: ", result)
			return result
		}
		i++
	}
	return result
}

func (caller *CallData) call(address string, input []byte, relatedAddress [][]byte) string {
	account := <-caller.server.availableAccounts
	log.Info("Use Account: ", account)
	conn, chCallData := caller.initCallConnection(account.Address)
	defer conn.GetTcpConnection().Close()
	hash := caller.sendCallData(conn, enterAddress(address), input, chCallData, account.Private,relatedAddress)
	for {

		select {
		case receiver := <-chCallData:
			log.Info("Hash on server", common.BytesToHash(hash.([]byte)))
			log.Info("Hash from chain", (receiver).(network.Receipt).Hash)
			if (receiver).(network.Receipt).Hash != common.BytesToHash(hash.([]byte)) {
				continue
			}
			go caller.server.GiveBackAccount(account)
			return (receiver).(network.Receipt).Value
		case <-time.After(10 * time.Second):
			go caller.server.GiveBackAccount(account)
			return "TimeOut"
		}
	}

}

//61.28.238.235:3011
type QueueLock struct {
	sync.Mutex
	queue map[string]bool
}

var queue = QueueLock{queue: make(map[string]bool)}


func (subscribe *Subscribe) subscribeChain(connRoot *cn.Connection) {
	contractsub := subscribe.server.contractABI
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contractsub["cake"].Address))
	log.Info("Listen address: ", contractsub["cake"].Address)
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contractsub["token"].Address))
	log.Info("Listen address: ", contractsub["token"].Address)
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contractsub["syrup"].Address))
	log.Info("Listen address: ", contractsub["syrup"].Address)
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contractsub["mtcv2"].Address))
	log.Info("Listen address: ", contractsub["mtcv2"].Address)
	// cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contractsub["cakepool"].Address))
	// log.Info("Listen address: ", contractsub["cakepool"].Address)
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contractsub["dummyToken"].Address))
	log.Info("Listen address: ", contractsub["dummyToken"].Address)
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contracts[0].Address))
	log.Info("Listen address: ", contracts[0].Address)
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contracts[1].Address))
	log.Info("Listen address: ", contracts[1].Address)
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contracts[2].Address))
	log.Info("Listen address: ", contracts[2].Address)
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contracts[3].Address))
	log.Info("Listen address: ", contracts[3].Address)
	cn.SendBytes(subscribe.server.config, connRoot, messages.SubscribeToAddress, common.FromHex(contracts[4].Address))
	log.Info("Listen address: ", contracts[4].Address)
}
func getPairConstant() Pair {
	return Pair{
		Address: "374b073c39741e8022ea217c2e4e46e0176413e5",
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
		Address: "198f1a030ddd3bbfCcfCCdED47E87eb7C30Fec81",
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
		Address: "FB1029ecda857c500Cc188437518ca9994C85452",
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
		Address: "FB725aaf7F11a3b7f192182fa1BF0755Ca993De0",
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

func (subscribe *Subscribe) handleSubscribeMessage() {
	server := subscribe.server
	broadcast := subscribe.subscribeChan
	cake := server.contractABI["cake"]
	token := server.contractABI["token"]
	dummy := server.contractABI["dummyToken"]
	syrup := server.contractABI["syrup"]
	mtcv2 := server.contractABI["mtcv2"]
	var caller CallData
	contract := server.contractABI["token0"]
	for conn, _ := range server.clients.data {
		caller = subscribe.server.clients.data[conn].caller
	}

	for {
		fmt.Println("start handleSubscribe")
		// capture event from chain
		msg := (<-broadcast).(network.EventI)
		// handle format event
		sendData := Message{}
		switch msg.Address {
		case strings.ToLower(mtcv2.Address):
			mtcv2.handleMtcv2Message(msg, &sendData,&caller)
		case strings.ToLower(cake.Address):
			cake.handleCakeMessage(msg, &sendData)
		case strings.ToLower(token.Address):
			token.handleTokenMessage(msg, &sendData)
		case strings.ToLower(dummy.Address):
			dummy.handleDummyMessage(msg, &sendData)
		case strings.ToLower(syrup.Address):
			syrup.handleSyrupMessage(msg, &sendData)
		case strings.ToLower(RouterContract.Address):
			handleRouterMessage(msg, &sendData)
		case strings.ToLower(FactoryContract.Address):
			handleFactoryMessage(msg, &sendData)
		case strings.ToLower(PairContract.Address):
			handlePairMessage(msg, &sendData)
		default:
			contract.handleAllTokenMessage(msg, &sendData)
		}

		//   Send it out to every client that is currently connected
		// for client := range clients {
		// 	sendQueue[client] <- sendData
		// }
		//   Send it out to every client that is currently connected
		log.Info(" - Send to all player - ")
		for _, client := range server.clients.data {
			client.sendChan <- sendData
			// client.caller.getFarmPoolInfoUpdate(msg.Data, msg.Topics)
			// conn := subscribe.connRoot
		}

	}
}
func readBool(word []byte) bool {
	for _, b := range word[:31] {
		if b != 0 {
			return false
		}
	}
	switch word[31] {
	case 0:
		return false
	case 1:
		return true
	default:
		return false
	}
}
func (contract *ContractABI) formatTopics(name string, result map[string]interface{}, topics map[int]interface{}) {
	i := 1

	for _, arg := range contract.Abi.Events[name].Inputs {
		if arg.Indexed {
			bytes := common.FromHex(topics[i].(string))
			switch arg.Type.T {
			case IntTy, UintTy:
				result[arg.Name] = ReadInteger(arg.Type, bytes)
			case BoolTy:
				result[arg.Name] = readBool(bytes)
			case AddressTy:
				result[arg.Name] = common.BytesToAddress(bytes)
			}
			i++
		}
	}
}

func (contract *ContractABI) decodeEvent(name string, data string, topics map[int]interface{}) map[string]interface{} {
	result := contract.decode(name, data).(map[string]interface{})
	contract.formatTopics(name, result, topics)
	return result
}
type Event2 struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func (contract *ContractABI) handleCakeMessage(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case contract.Abi.Events["Transfer"].ID.String()[2:]:
		name := contract.Abi.Events["Transfer"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	case contract.Abi.Events["Approval"].ID.String()[2:]:
		name := contract.Abi.Events["Approval"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	default:
	}
}
func (contract *ContractABI) handleSyrupMessage(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case contract.Abi.Events["Transfer"].ID.String()[2:]:
		name := contract.Abi.Events["Transfer"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	case contract.Abi.Events["Approval"].ID.String()[2:]:
		name := contract.Abi.Events["Approval"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	default:
	}
}
func (contract *ContractABI) handleDummyMessage(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case contract.Abi.Events["Transfer"].ID.String()[2:]:
		name := contract.Abi.Events["Transfer"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	case contract.Abi.Events["Approval"].ID.String()[2:]:
		name := contract.Abi.Events["Approval"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	default:
	}
}
func (contract *ContractABI) handleTokenMessage(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case contract.Abi.Events["Transfer"].ID.String()[2:]:
		name := contract.Abi.Events["Transfer"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	case contract.Abi.Events["Approval"].ID.String()[2:]:
		name := contract.Abi.Events["Approval"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	default:
	}
}
func (contract *ContractABI) handleAllTokenMessage(msg network.EventI, sendData *Message) {

	switch msg.Event {
	case contract.Abi.Events["Transfer"].ID.String()[2:]:
		name := contract.Abi.Events["Transfer"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	case contract.Abi.Events["Approval"].ID.String()[2:]:
		name := contract.Abi.Events["Approval"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	default:
	}
}
func (caller *CallData) getFarmPoolInfoUpdate(data ...interface{}) {

	//tìm kiếm token được deposit/withraw
	var pool Pool
	var list ListUpdate
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"poolid": formatNumber(data[1].(map[int]interface{})[2].(string))}
	err := collection.FindOne(ctx, filter).Decode(&pool)
	if err != nil {
		fmt.Println("Find database failled")
		return
	}
	//lấy số liệu APR, multiplier, liquidity của token đó
	id:= pool.ID
	address := pool.Address
	multiplier,_ := new(big.Int).SetString(pool.Multiplier,10)
	apr,_ := new(big.Int).SetString(pool.APR,10)
	liquidity,_ := new(big.Int).SetString(pool.Liquidity,10)
	amount,_ := new(big.Int).SetString(formatNumber(data[0].(string)),10)
	var newLiquidity *big.Int
	switch data[2]{
		case "Deposit":
			newLiquidity = new(big.Int).Add(liquidity,amount)
		case "Withdraw":
			newLiquidity = new(big.Int).Sub(liquidity,amount)
	}
	
	fmt.Println("multiplier:",multiplier)
	fmt.Println("apr:",apr)
	fmt.Println("liquidity:",liquidity)
	fmt.Println("newLiquidity:",newLiquidity)
	fmt.Println("id:",id)


	//Tính lại database mới của token đó
		var APR *big.Int
		cakePerblock := 1
		priceOfCake := big.NewInt(int64(20)) //getAmountOut in router????
		cakePerYear := big.NewInt(int64(cakePerblock * 60 * 60 * 24 * 365))
		volume24h := 50000 // dang fix cung theo usd
		amountOfCakePerYear := big.NewInt(0).Mul(multiplier, cakePerYear)
		totalValueOfCake := big.NewInt(0).Mul(amountOfCakePerYear, priceOfCake)
		if newLiquidity.Cmp(big.NewInt(0)) == 0 {

			APR = big.NewInt(0)
			fmt.Println("Liquidity is zero")
		} else {
			farmBaseReward := big.NewInt(0).Div(big.NewInt(0).Mul(totalValueOfCake, new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)), newLiquidity) //Percentage
			//farmBaseReward= totalValueOfCake)/newLiquidity*100/10^18
			totalFee := volume24h * 17 / 10000 * 365
			lpReward := big.NewInt(0).Div(big.NewInt(int64(totalFee)), big.NewInt(0).Div(newLiquidity, new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)))
			// lpReward = totalFee/(newLiquidity/10^18)*100
			fmt.Println("farmBaseReward:", farmBaseReward)
			fmt.Println("lpReward:", lpReward)
			APR = big.NewInt(0).Add(farmBaseReward, lpReward)
		}
		log.Info("GetAprPool - Result - ", APR)
		newPool := Pool{
			ID:         id,
			Address:    address,
			APR:        APR.String(),
			Multiplier: multiplier.String(),
			Liquidity:  newLiquidity.String(),
		}
		//Update vào database
		ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		result2,err := collection.UpdateOne(ctx1,Pool{ID: id}, bson.M{"$set": newPool})
		if err != nil {
			fmt.Println("Update database failled")
			return
		}else{
			fmt.Println("UpdateOne() result:", result2)
		}
		out, err := json.Marshal(newPool)
		if err != nil {
			panic (err)
		}
		result := string(out)
		log.Info("getFarmPoolInfoUpdate- Result - ", result)
		fmt.Printf("Type out:%T",string(out))
		list.Lock()
		list.data = append(list.data, newPool)
		list.Unlock()
	go sentToClient1("getFarmPoolInfoUpdate", list.data)
}

func (contract *ContractABI) handleMtcv2Message(msg network.EventI, sendData *Message,caller *CallData) {
	switch msg.Event {
	case contract.Abi.Events["Deposit"].ID.String()[2:]:
		name := contract.Abi.Events["Deposit"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
		// log.Info(" Deposit Event:",msg.Data)
		// for _, client := range server.clients.data {
		// 	client.sendChan <- sendData
		// }
		// chFarm <- msg.Data
		// {From:0EE5f1926C29B6e02d151A31A31cabBF15001b32 {Deposit map[amount:100000000000000000000 pid:1 user:0x95d5e615BfaD47B785723a1dCCd1Ee7a4E8c7580]}} 
		caller.getFarmPoolInfoUpdate(msg.Data, msg.Topics, name)
	case contract.Abi.Events["UpdatePool"].ID.String()[2:]:
		name := contract.Abi.Events["UpdatePool"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
	case contract.Abi.Events["Withdraw"].ID.String()[2:]:
		name := contract.Abi.Events["Withdraw"].Name
		*sendData = Message{
			"From:" + contract.Address,
			Event2{
				name,
				contract.decodeEvent(name, msg.Data, msg.Topics),
			},
		}
		caller.getFarmPoolInfoUpdate(msg.Data, msg.Topics, name)
	default:
	}
}

func handleRouterMessage(msg network.EventI, sendData *Message) {
	switch msg.Event {
	case RouterContract.AddLiquidity.Hash:
		*sendData = Message{
			"From:" + RouterContract.Address,
			Event2{RouterContract.AddLiquidity.Name, formatEvent(
				string(msg.Data),
				RouterContract.AddLiquidity.Format,
			)},
		}
	case RouterContract.AddLiquidity.Hash:
		*sendData = Message{
			"From:" + RouterContract.Address,
			Event2{RouterContract.GetPriceList.Name, formatEvent(
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
			Event2{FactoryContract.PairCreated.Name, formatEvent(
				string(msg.Data),
				FactoryContract.PairCreated.Format,
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
			Event2{PairContract.Transfer.Name, formatEvent(
				string(msg.Data),
				PairContract.Transfer.Format,
			)},
		}
	case PairContract.Sync.Hash:
		*sendData = Message{
			"From:" + PairContract.Address,
			Event2{PairContract.Sync.Name, formatEvent(
				string(msg.Data),
				PairContract.Sync.Format,
			)},
		}
	case PairContract.Mint.Hash:
		*sendData = Message{
			"From:" + PairContract.Address,
			Event2{PairContract.Mint.Name, formatEvent(
				string(msg.Data),
				PairContract.Mint.Format,
			)},
		}
	case PairContract.Burn.Hash:
		*sendData = Message{
			"From:" + PairContract.Address,
			Event2{PairContract.Burn.Name, formatEvent(
				string(msg.Data),
				PairContract.Burn.Format,
			)},
		}

	case PairContract.Approval.Hash:
		*sendData = Message{
			"From:" + PairContract.Address,
			Event2{PairContract.Approval.Name, formatEvent(
				string(msg.Data),
				PairContract.Approval.Format,
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
func formatData(dataType string, data interface{}) interface{} {
	switch dataType {
	case "string":
		return data.(string)
	case "bool":
		return data.(bool)
	case "address":
		return common.HexToAddress(data.(string))
	case "uint8":
		intVar, err := strconv.Atoi(data.(string))
		if err != nil {
			log.Warn("Conver Uint8 fail", err)
			return nil
		}
		return uint8(intVar)
	case "uint", "uint256":
		nubmer := big.NewInt(0)
		nubmer, ok := nubmer.SetString(data.(string), 10)
		if !ok {
			log.Warn("Format big int: error")
			return nil
		}
		return nubmer
	// case "array","slice" :
	// rv := reflect.ValueOf(data)
	// var out []interface{}
	// 	for i := 0; i < rv.Len(); i++ {
	// 		out = append(out, rv.Index(i).Interface())
	// 	}
	// 		tk0 := common.HexToAddress(data[1].(string))
	// tk1 := common.HexToAddress(data[2].(string))
	// path := []common.Address{tk0, tk1}

	// return data
	default:
		return data
	}
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

func (caller *CallData) sendCallData(
	// config config.Config,
	conn *cn.Connection,
	toAddress common.Address,
	input []byte,
	chCallData chan interface{},
	private string,
	relatedAddress [][]byte,
	)interface{} {
	amount := uint256.NewInt(0)
	transferFee := uint256.NewInt(1)
	// Get info from secretkey
	bRootPrikey, bRootPubkey, bRootAddress := core_crypto.GenerateKeyPairFromSecretKey(private)
	rootAddress := common.BytesToAddress(bRootAddress)
	sendGetAccountState(caller.server.config, conn, rootAddress)

	var accountInfo *pb.AccountState
	select {
	case temp := <-chCallData:
		result, ok := temp.(*pb.AccountState)
		if !ok {
			return nil
		}
		accountInfo = result
	case <-time.After(5 * time.Second):
		return nil
	}
	// Get last transaction
	// fmt.Println("accountInfo:", accountInfo)

	lastTransaction := cc.GetEmptyTransaction()

	// if hex.EncodeToString(lastTransaction.Hash) != hex.EncodeToString(accountInfo.LastHash) {
	// 	lastTransaction = readDataLastHash(hex.EncodeToString(accountInfo.LastHash))
	// }
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
		RelatedAddresses: relatedAddress,
	}
	// smartContractData.CallData.FeeType = pb.FEETYPE_SMART_CONTRACT_CHARGE_FEE 

	txt := &pb.Transaction{
		FromAddress:         bRootAddress,
		ToAddress:           toAddress.Bytes(),
		PubKey:              bRootPubkey,
		PendingUse:          u256PendingUse.Bytes(),
		Balance:             balance.Bytes(),
		Amount:              amount.Bytes(),
		Fee:                 transferFee.Bytes(),
		Data:                smartContractData,
		LastHash:      accountInfo.LastHash,
		LastDeviceKey: common.FromHex("0000000000000000000000000000000000000000000000000000000000000000"),
		NewDeviceKey:  common.FromHex("290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"),
	}
	hash := cc.GetTransactionHash(txt)
	txt.Hash = hash
	// Sign and VerifyTransaction
	txt.Sign = core_crypto.Sign(bRootPrikey, hash)

	txtB, err := proto.Marshal(txt)
	if err != nil {
		log.Fatal(err)
	}

	// err = os.WriteFile(fmt.Sprintf("./datas/%s", hex.EncodeToString(txt.Hash)), txtB, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cn.SendBytes(caller.server.config, conn, messages.SendTransaction, txtB)

	transactionsDB.GetInstanceTransactionsDB().PendingTransaction = txt
	fmt.Printf("Transaction sent")
	return hash
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
	config config.Config,
	parentConn *cn.Connection,
	address common.Address,
) {
	cn.SendBytes(config, parentConn, messages.GetAccountState, address.Bytes())
	fmt.Println("Sended Account Info")
}
func enterAddress(message string) common.Address {
	address := strings.Replace(message, "\n", "", -1)
	address = strings.ToLower(message)
	return common.HexToAddress(address)
}

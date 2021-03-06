package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/YushinJung/NomadCoin/blockchain"
	"github.com/YushinJung/NomadCoin/p2p"
	"github.com/YushinJung/NomadCoin/utils"
	"github.com/YushinJung/NomadCoin/wallet"
	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayLoad struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type addPeerPayLoad struct {
	Address, Port string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See all Block",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
		{
			URL:         url("/mempool"),
			Method:      "Get",
			Description: "Get Mempool",
		},
		{
			URL:         url("/transaction"),
			Method:      "POST",
			Description: "Add a transaction",
			Payload:     "data:addTxPayLoad",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to Web Sockets",
		},
	}
	//rw.Header().Add("Content-Type", "application/json")
	// middleware ????????? ?????? ?????????.
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//rw.Header().Add("Content-Type", "application/json")
		//middleware ????????? ?????? ?????????.
		// json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain()))
	case "POST":
		newBlock := blockchain.Blockchain().AddBlock()
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated) // header??? created ????????? ???????????? ???
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // map[id: ]
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	ecoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		ecoder.Encode(errorResponse{ErrorMessage: fmt.Sprint(err)})
	} else {
		ecoder.Encode(block)
	}
}

func jsonConentTypeMiddelWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// do stuff here
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
		// point going next
	})
	// handler ??? interface ??????. ??? interface??? ServeHTTP???
}

func loggerMiddelware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// do stuff here
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
		// point going next
	})
	// handler ??? interface ??????. ??? interface??? ServeHTTP???
}

func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), rw)
	// blockchain????????? ???????????? ??? ?????? ?????? ?????? ????????? ?????? ????????? ??????????????? ?????? ???
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // get variable
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true": // total amount needs to get extra struct since it was not defined before
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(rw).Encode(balanceResponse{
			Address: address,
			Balance: amount,
		})
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	blockchain.StatusMempool(rw)
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayLoad
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload)) // json??? ????????? addTxPayLoad struct??? ??????????????? ??????
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		// if this happens we have to kill this function
		return
	}
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(struct {
		Address string `json:"address"`
	}{Address: address})
	//json.NewEncoder(rw).Encode(myWalletResponse{Address: address})

}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayLoad
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port, true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonConentTypeMiddelWare, loggerMiddelware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance)
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("POST", "GET")
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/YushinJung/NomadCoin/blockchain"
	"github.com/YushinJung/NomadCoin/utils"
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
type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
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
	}
	//rw.Header().Add("Content-Type", "application/json")
	// middleware 추가로 필요 없어짐.
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//rw.Header().Add("Content-Type", "application/json")
		//middleware 추가로 필요 없어짐.
		// json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		var aBB addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&aBB))
		// r.Body 에서 받아와서 addBlockBody 에 넣을 것
		blockchain.Blockchain().AddBlock(aBB.Message)
		rw.WriteHeader(http.StatusCreated) // header로 created 됐다고 알려주는 것
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
	// handler 는 interface 이다. 이 interface는 ServeHTTP를
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonConentTypeMiddelWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

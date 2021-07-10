package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/YushinJung/NomadCoin/blockchain"
)

const port string = ":4000"

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	// to write somewhere (small)
	// read data (big)

	templ := template.Must(template.ParseFiles("templates/home.gohtml"))
	// templ, err := template.ParseFiles("templates/home.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templ.Execute(rw, data)
}
func main() {
	// http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", home)
	fmt.Printf("listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil)) // log.Fatal will stop if there is error from input
}

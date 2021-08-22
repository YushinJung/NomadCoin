package p2p

import (
	"fmt"
	"net/http"

	"github.com/YushinJung/NomadCoin/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	for {
		fmt.Println("Waiting 4 message ...")
		_, p, err := conn.ReadMessage() // blocking function
		fmt.Println("Message arrived")
		// websocket connection
		utils.HandleErr(err)
		fmt.Printf("%s", p)
	}
}

package p2p

import (
	"fmt"
	"net/http"
	"time"

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
		if err != nil {
			break
		}
		fmt.Println("Message arrived")
		// websocket connection
		utils.HandleErr(err)
		fmt.Printf("Just got: %s\n", p)
		time.Sleep(2 * time.Second)
		message := fmt.Sprintf("New Message: %s", p)
		utils.HandleErr(conn.WriteMessage(websocket.TextMessage, []byte(message)))
	}
}

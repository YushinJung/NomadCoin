package p2p

import (
	"fmt"
	"net/http"

	"github.com/YushinJung/NomadCoin/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Port :3000 will upgrade the request from :4000
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	openPort := r.URL.Query().Get("openPort")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)

}

func AddPeer(address, port, openPort string) {
	formatedPort := utils.Splitter(openPort, ":", 1)
	// from :4000 is requesting an upgrade at the port :3000
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, formatedPort), nil)
	// 원래 nil 이 들어가는 부분에 requestheader들어가서 authenticate을 하는데,
	// 여기서는 그냥 nil을 쓰자.
	utils.HandleErr(err)
	initPeer(conn, address, port)
}

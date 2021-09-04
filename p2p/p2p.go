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
	fmt.Printf("%s wants an upgrade\n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)

	// conn 에 바로 message를 보내면, peer를 생성할 때만 message를 쓸 수 있음.
	// conn이 들어 있는 peer는 외부 변수로 peer의 inbox에 message를 보내고
	// inbox(channel)에 message가 들어올 때, write기능 발생 시킴.
}

func AddPeer(address, port, openPort string) {
	fmt.Printf("%s wants to connect to port %s\n", openPort, port)
	formatedPort := utils.Splitter(openPort, ":", 1)
	// from :4000 is requesting an upgrade at the port :3000
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, formatedPort), nil)
	// 원래 nil 이 들어가는 부분에 requestheader들어가서 authenticate을 하는데,
	// 여기서는 그냥 nil을 쓰자.
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	sendNewestBlock(peer)
}

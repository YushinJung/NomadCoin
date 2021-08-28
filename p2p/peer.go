package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	conn    *websocket.Conn
	inbox   chan []byte
	key     string
	address string
	port    string
}

func (p *peer) close() {
	p.conn.Close()
	delete(Peers, p.key) // 스스로 지우기가 가능하군
}

func (p *peer) read() {
	// delete peer in case of error
	defer p.close() // run after this function finishes
	for {
		_, m, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

func (p *peer) write() {

	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			break
		}
		err := p.conn.WriteMessage(websocket.TextMessage, m)
		if err != nil {
			break
		}
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		address: address,
		key:     key,
		port:    port,
	}
	go p.read() // pear 생성 후 read를 계속 진행 시킬 것
	go p.write()
	Peers[key] = p
	return p
}

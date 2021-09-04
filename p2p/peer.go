package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type peers struct {
	v map[string]*peer
	m sync.Mutex
}

var Peers peers = peers{
	v: make(map[string]*peer),
}

type peer struct {
	conn    *websocket.Conn
	inbox   chan []byte
	key     string
	address string
	port    string
}

func AllPeers(p *peers) []string {
	p.m.Lock()
	defer p.m.Unlock()

	var keys []string
	for key := range p.v {
		keys = append(keys, key)
	}
	return keys
}

func (p *peer) close() {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
	delete(Peers.v, p.key) // 스스로 지우기가 가능하군
}

func (p *peer) read() {
	// delete peer in case of error
	defer p.close() // run after this function finishes
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m) // connection에 message가 올때까지 기다린다
		if err != nil {
			break
		}
		hanldeMsg(&m, p)
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
	Peers.m.Lock()
	defer Peers.m.Unlock()
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
	Peers.v[key] = p
	return p
}

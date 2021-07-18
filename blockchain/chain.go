package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/YushinJung/NomadCoin/db"
	"github.com/YushinJung/NomadCoin/utils"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain // will not be shared
var once sync.Once

func (b *blockchain) restore(data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	utils.HandleErr(decoder.Decode(b))
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	// by adding block we should update newesthash and height
	b.NewestHash = block.Hash
	b.Height = block.Height
	// need to update db
	b.persist()
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			fmt.Printf("NewestHash: %s\nHeight:%d\n", b.NewestHash, b.Height)
			//search for checkpoint in db
			checkpoint := db.Blockchain()
			if checkpoint == nil {
				// nothing is at the block chain
				// height will be 0
				b.AddBlock("Genesis")
			} else {
				// restore b from bytes
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}
		})
	}
	fmt.Printf("NewestHash: %s\nHeight:%d\n", b.NewestHash, b.Height)
	return b
}

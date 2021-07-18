package blockchain

import (
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
	utils.FromBytes(b, data)
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

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}
	return blocks
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			//search for checkpoint in db
			checkpoint := db.Blockchain()
			if checkpoint == nil {
				// nothing is at the block chain
				// height will be 0
				b.AddBlock("Genesis")
			} else {
				// restore b from bytes
				b.restore(checkpoint)
			}
		})
	}
	return b
}

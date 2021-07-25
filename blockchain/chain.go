package blockchain

import (
	"sync"

	"github.com/YushinJung/NomadCoin/db"
	"github.com/YushinJung/NomadCoin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5 // change difficulty every  5 block interval
	blockInterval      int = 2 // we want block every two minutes
	allowedRangeTime   int = 3 // allowed range of time to change difficulty
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDificulty"`
}

var b *blockchain // will not be shared
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	// by adding block we should update newesthash and height
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
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
	// changed order backward
	// for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
	// 	blocks[i], blocks[j] = blocks[j], blocks[i]
	// }
	return blocks
}

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[-1+difficultyInterval]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRangeTime) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRangeTime) {
		return b.CurrentDifficulty - 1
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 { // first block
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return b.recalculateDifficulty()
	} else {
		// previous difficulty
		return b.CurrentDifficulty
	}
}

func (b *blockchain) txOuts() []*TxOut {
	// get all transaction outputs
	var txOuts []*TxOut
	blocks := b.Blocks()
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
			// TxOuts is also a slice fo TxOut
		}
	}
	return txOuts
}

func (b *blockchain) TxOutsByAddress(address string) []*TxOut {
	var ownedTxOuts []*TxOut
	txOuts := b.txOuts()
	for _, txOut := range txOuts {
		if txOut.Owner == address {
			ownedTxOuts = append(ownedTxOuts, txOut)
		}
	}
	return ownedTxOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	txOuts := b.TxOutsByAddress(address)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			//search for checkpoint in db
			checkpoint := db.Blockchain()
			if checkpoint == nil {
				// nothing is at the block chain
				// height will be 0
				b.AddBlock()
			} else {
				// restore b from bytes
				b.restore(checkpoint)
			}
		})
	}
	return b
}

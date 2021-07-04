package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	// save data in the blcok only
	data     string
	hash     string // connected
	prevHash string
}

// hash function is one-way function
// if the data of hash is changed in the middle
// all of the hash will be changed after that hash

func main() {
	genesisBlock := block{"Genesis Block", "", ""}
	hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.prevHash))
	genesisBlock.hash = fmt.Sprintf("%x", hash)
	fmt.Println(genesisBlock)
}

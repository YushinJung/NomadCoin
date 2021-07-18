package main

import (
	"github.com/YushinJung/NomadCoin/blockchain"
)

func main() {
	blockchain.Blockchain().AddBlock("First")
	blockchain.Blockchain().AddBlock("Second")
	blockchain.Blockchain().AddBlock("Third")
}

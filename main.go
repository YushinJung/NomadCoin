package main

import (
	"github.com/YushinJung/NomadCoin/blockchain"
	"github.com/YushinJung/NomadCoin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
}

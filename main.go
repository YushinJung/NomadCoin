package main

import (
	"github.com/YushinJung/NomadCoin/cli"
	"github.com/YushinJung/NomadCoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}

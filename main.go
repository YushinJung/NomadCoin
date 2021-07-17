package main

import (
	"github.com/YushinJung/NomadCoin/explorer"
	"github.com/YushinJung/NomadCoin/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}

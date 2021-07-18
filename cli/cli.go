package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/YushinJung/NomadCoin/explorer"
	"github.com/YushinJung/NomadCoin/rest"
)

func usage() {
	fmt.Printf("Welcome to Yushin Coin\n\n")
	fmt.Printf("Pelase use the following flags:\n\n")
	fmt.Printf("-port=4000: 	Set the PORT of the server\n")
	fmt.Printf("-mode=rest		'html' vs 'rest'\n\n")
	os.Exit(0)
}

func Start() {
	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "'html' vs 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		// start with rest api
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}

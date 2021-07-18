package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("Welcome to Yushin Coin\n\n")
	fmt.Printf("Pelase use the following commands:\n\n")
	fmt.Printf("explorer: 	Start the HTML Explore\n")
	fmt.Printf("rest:		Start the REST API (recommended)\n\n")
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	rest := flag.NewFlagSet("rest", flag.ExitOnError)
	portFlag := rest.Int("port", 4000, "Sets the port of the server")

	switch os.Args[1] {
	case "explorer":
		fmt.Println("Start Explorer")
	case "rest":
		rest.Parse(os.Args[2:])
		fmt.Println("Start REST API")
	default:
		usage()
	}
	fmt.Println(*portFlag)
}

package main

import (
	"fmt"

	"github.com/YushinJung/NomadCoin/person"
)

func main() {
	yushin := person.Person{}
	yushin.SetDetails("yushin", 12)
	fmt.Println("Main", yushin)
}

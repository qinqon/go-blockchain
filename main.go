package main

import (
	"github.com/qinqon/go-blockchain/blockchain"
	"log"
	"os"
)

func main() {
	address := ":5000"
	if len(os.Args) >= 2 {
		address = os.Args[1]
	}
	node := blockchain.NewNode(address)

	if err := node.Start(address); err != nil {
		log.Fatal(err)
	}
}

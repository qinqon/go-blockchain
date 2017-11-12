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

	log.Printf("Starting blockchain node at %v", address)
	if err := node.HttpServer().ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

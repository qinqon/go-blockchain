package main

import (
	"github.com/qinqon/go-blockchain/blockchain"
)

func main() {
	node := blockchain.NewNode()
	node.Start(":5000")
}

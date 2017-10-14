package main

import ("github.com/qinqon/go-blockchain/blockchain")

func main() {
   bc := blockchain.New()
   bc.StartWebServer()
}

package blockchain

import "testing"

func TestBlockchain(t *testing.T) {
   blockchain := Blockchain{[]Blocks, []Transaction}
   print(blockchain)
}

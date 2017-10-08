package blockchain

import "time"
import "log"
import "encoding/json"
import "crypto/sha256"
import "encoding/hex"

type Blockchain struct{
   chain []*Block
   currentTransactions []Transaction
}
func NewBlockchain() *Blockchain {
   bc := Blockchain{}
   return &bc
}

func Hash(block *Block) string{
   marshalledBlock, err := json.Marshal(*block)
   if err != nil {
      log.Fatal(err)
   }
   hashedBlock := sha256.Sum256(marshalledBlock)
   hexBlock := hex.EncodeToString(hashedBlock[:])
   return hexBlock
}

func (b *Blockchain) NewBlock(proof int, previousHash int) *Block{
   block := Block{
      Index: len(b.chain) + 1,
      Timestamp: time.Time{},
      Transactions: b.currentTransactions,
      Proof: proof,
      PreviousHash: previousHash}
   b.currentTransactions = []Transaction{}
   b.chain = append(b.chain, &block)
   return &block;
}

func (b *Blockchain) NewTransaction(sender string, recipient string, amount int) int{
   b.currentTransactions = append(b.currentTransactions, Transaction{sender, recipient, amount})
   return b.LastBlock().Index + 1
}


func (b *Blockchain) LastBlock() *Block{
   return b.chain[len(b.chain) - 1]
}

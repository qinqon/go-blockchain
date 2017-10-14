package blockchain

import "time"
import "log"
import "encoding/json"
import "crypto/sha256"
import "encoding/hex"
import "fmt"

type Blockchain struct{
   chain []*Block
   currentTransactions []Transaction
}
func New() *Blockchain {
   bc := Blockchain{}
   return &bc
}

func Hash(block *Block) string{
   //TODO: Sort dictionary
   marshalledBlock, err := json.Marshal(*block)
   if err != nil {
      log.Fatal(err)
   }
   hashedBlock := sha256.Sum256(marshalledBlock)
   hexBlock := hex.EncodeToString(hashedBlock[:])
   return hexBlock
}

func (b *Blockchain) ValidProof(lastProof, proof int) bool {
   guess := []byte(fmt.Sprintf("%i%i", lastProof, proof))
   guessHash := sha256.Sum256(guess)
   //FIXME: Don't convert to string, compare the slice
   return string(guessHash[len(guessHash)-4]) == "0000"
}


func (b *Blockchain) ProofOfWork(lastProof int) int {
   proof := 0
   for ; b.ValidProof(proof, lastProof); proof++ {
   }
   return proof
}

func (b *Blockchain) NewBlock(proof int, previousHash string) *Block{
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

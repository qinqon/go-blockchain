package blockchain

import (
   "errors"
   "time"
   "log"
   "encoding/json"
   "crypto/sha256"
   "encoding/hex"
   "fmt"
   "net/url"
)

type Blockchain struct{
   nodes map[*url.URL]bool
   chain []*Block
   currentTransactions []Transaction
}
func New() *Blockchain {
   bc := Blockchain{}
   // Genesis block
   bc.currentTransactions = []Transaction{}
   bc.NewBlockWithPreviousHash(100, "1")
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

func (b Blockchain) Chain() []*Block {
   return b.chain
}

func ValidProof(lastProof, proof int) bool {
   guess := []byte(fmt.Sprintf("%i%i", lastProof, proof))
   guessHash := sha256.Sum256(guess)
   guessHashEncoded := hex.EncodeToString(guessHash[:])
   return guessHashEncoded[len(guessHashEncoded)-4:len(guessHashEncoded)] == "0000"
}


func ProofOfWork(lastProof int) int {
   proof := 0
   for ValidProof(proof, lastProof) == false {
      proof++
   }
   return proof
}


func (b *Blockchain) NewBlock(proof int) *Block {
   return b.NewBlockWithPreviousHash(proof, Hash(b.chain[len(b.chain)-1]))
}

func (b *Blockchain) NewBlockWithPreviousHash(proof int, previousHash string) *Block{
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

func (b *Blockchain) NewTransaction(transaction Transaction) int{
   b.currentTransactions = append(b.currentTransactions, transaction)
   return b.LastBlock().Index + 1
}


func (b *Blockchain) LastBlock() *Block{
   return b.chain[len(b.chain) - 1]
}

func (b *Blockchain) RegisterNode(address string) {
   u, err := url.Parse(address)
   if err  != nil {
      log.Fatal(err)
   }
   b.nodes[u] = true
}

func ValidChain(chain []*Block) error {
   if len(chain) <= 1 {
      return nil
   }
   var lastBlock *Block = chain[0]
   for _, block := range chain[1:len(chain)-1]{
      if ! ValidProof(lastBlock.Proof, block.Proof) {
         return errors.New("Invalid proof")
      }
      if block.PreviousHash != Hash(lastBlock) {
         return errors.New("Invalid privious hash")
      }
      lastBlock = block
   }
   return nil
}

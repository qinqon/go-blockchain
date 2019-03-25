package blockchain

import (
	"log"
	"time"
)

type Blockchain struct {
	chain               []*Block
	currentTransactions []Transaction
}

func NewBlockchain() *Blockchain {
	bc := Blockchain{}
	// Genesis block
	bc.currentTransactions = []Transaction{}
	bc.NewBlockWithPreviousHash(100, "1")
	return &bc
}

func (b Blockchain) Chain() []*Block {
	return b.chain
}

func (b *Blockchain) NewBlock(proof int) *Block {
	previousBlockHash, err := b.chain[len(b.chain)-1].CalculateHash()
	if err != nil {
		log.Fatal(err)
	}
	return b.NewBlockWithPreviousHash(proof, previousBlockHash)
}

func (b *Blockchain) NewBlockWithPreviousHash(proof int, previousHash string) *Block {
	block := Block{
		Index:        len(b.chain) + 1,
		Timestamp:    time.Time{},
		Transactions: b.currentTransactions,
		Proof:        proof,
		PreviousHash: previousHash}
	b.currentTransactions = []Transaction{}
	b.chain = append(b.chain, &block)
	return &block
}

func (b *Blockchain) NewTransaction(transaction Transaction) int {
	b.currentTransactions = append(b.currentTransactions, transaction)
	return b.LastBlock().Index + 1
}

func (b *Blockchain) LastBlock() *Block {
	return b.chain[len(b.chain)-1]
}

func ValidateChain(chain []*Block) error {
	if len(chain) <= 1 {
		return nil
	}
	var previousBlock *Block = chain[0]
	for _, block := range chain[1:] {
		if err := block.Validate(previousBlock); err != nil {
			return err
		}
		previousBlock = block
	}
	return nil
}

func (b *Blockchain) Validate() error {
	return ValidateChain(b.Chain())
}

func (b *Blockchain) ResolveConflict(newChain []*Block) {
	if len(newChain) < len(b.chain) {
		return
	}
	if err := ValidateChain(newChain); err == nil {
		b.chain = newChain
	}
}

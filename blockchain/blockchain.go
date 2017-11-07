package blockchain

import (
	"log"
	"net/url"
	"time"
)

type Blockchain struct {
	nodes               map[*url.URL]bool
	chain               []*Block
	currentTransactions []Transaction
}

func New() *Blockchain {
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

func (b *Blockchain) Validate() error {
	if len(b.Chain()) <= 1 {
		return nil
	}
	var previousBlock *Block = b.Chain()[0]
	for _, block := range b.Chain()[1:] {
		if err := block.Validate(previousBlock); err != nil {
			return err
		}
		previousBlock = block
	}
	return nil
}

func (b *Blockchain) RegisterNode(address string) error {
	u, err := url.Parse(address)
	if err != nil {
		return err
	}
	b.nodes[u] = true
	return nil
}

func (b *Blockchain) ResolveConflict(newBlockchain *Blockchain) {
	if len(newBlockchain.Chain()) < len(b.chain) {
		return
	}
	if err := newBlockchain.Validate(); err == nil {
		b.chain = newBlockchain.chain
	}
}

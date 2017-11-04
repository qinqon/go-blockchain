package blockchain

import(
   "testing"
   "github.com/stretchr/testify/assert"
)

func TestEmptyChainIsValid(t *testing.T) {
   var emptyChain []*Block
   assert.True(t, ValidChain(emptyChain), "Invalid chain")
}

func TestChainWithJustGenesIsBlockIsValid(t *testing.T) {
   blockchain := New()
   assert.True(t, ValidChain(blockchain.Chain()), "Invalid chain")
}

func TestTwoBlocksChainWithCorrectProofIsValid(t *testing.T) {
   blockchain := New()
   lastProof := blockchain.LastBlock().Proof
   proof := ProofOfWork(lastProof)
   blockchain.NewBlock(proof)
   assert.True(t, ValidChain(blockchain.Chain()), "Invalid chain")
}

func TestTwoBlocksChainWithWrongProofIsInvalid(t *testing.T) {
   blockchain := New()
   lastProof := blockchain.LastBlock().Proof
   proof := ProofOfWork(lastProof + 1)
   blockchain.NewBlock(proof)
   assert.False(t, ValidChain(blockchain.Chain()), "Unexpected valid chain")
}

package blockchain

import(
   "testing"
   "github.com/stretchr/testify/assert"
)

func TestEmptyChainIsValid(t *testing.T) {
   var emptyChain []*Block
   assert.Nil(t, ValidChain(emptyChain), "Invalid chain")
}

func TestChainWithJustGenesIsBlockIsValid(t *testing.T) {
   blockchain := New()
   assert.Nil(t, ValidChain(blockchain.Chain()), "Invalid chain")
}

func TestTwoBlocksChainWithCorrectProofIsValid(t *testing.T) {
   blockchain := New()
   lastProof := blockchain.LastBlock().Proof
   proof := ProofOfWork(lastProof)
   blockchain.NewBlock(proof)
   assert.Nil(t, ValidChain(blockchain.Chain()), "Invalid chain")
}

func TestTwoBlocksChainWithWrongProofIsInvalid(t *testing.T) {
   blockchain := New()
   lastProof := blockchain.LastBlock().Proof
   proof := ProofOfWork(lastProof + 1)
   blockchain.NewBlock(proof)
   assert.NotNil(t, ValidChain(blockchain.Chain()), "Unexpected valid chain")
}

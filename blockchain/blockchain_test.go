package blockchain

import(
   "testing"
   "github.com/stretchr/testify/assert"
)

func TestInvalidProofReturnsFalse(t *testing.T) {
   lastProof := 33
   proof := 10033575
   assert.False(t, ValidProof(lastProof, proof), "33 and 66 as last and current proof is not valid")
}

func TestValidProofReturnsTrue(t *testing.T) {
   lastProof := 33
   proof := 10033575
   assert.True(t, ValidProof(lastProof, proof), "")
}

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

   // Bad proof
   blockchain.NewBlock(lastProof)

   assert.NotNil(t, ValidChain(blockchain.Chain()), "Unexpected valid chain")
}

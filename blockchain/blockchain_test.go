package blockchain

import(
   "testing"
   "github.com/stretchr/testify/assert"
)

func TestInvalidProofReturnsFalse(t *testing.T) {
   lastProof := 33
   proof := 66
   assert.Errorf(t, ValidProof(lastProof, proof), "last proof %v and current proof %v are not valid", lastProof, proof)
}

func TestValidProofReturnsTrue(t *testing.T) {
   lastProof := 33
   proof := 10033575
   assert.Errorf(t, ValidProof(lastProof, proof), "last proof %v and current proof %v are valid", lastProof, proof)
}

func TestChainWithJustGenesIsBlockIsValid(t *testing.T) {
   blockchain := New()
   assert.NoError(t, blockchain.Validate(), "Invalid chain")
}

func TestTwoBlocksChainWithCorrectProofIsValid(t *testing.T) {
   blockchain := New()
   lastProof := blockchain.LastBlock().Proof
   proof := ProofOfWork(lastProof)
   blockchain.NewBlock(proof)
   assert.NoError(t, blockchain.Validate(), "Invalid chain")
}

func TestTwoBlocksChainWithWrongProofIsInvalid(t *testing.T) {
   blockchain := New()
   lastProof := blockchain.LastBlock().Proof

   // Bad proof
   blockchain.NewBlock(lastProof)
   assert.NotNil(t, blockchain.Validate(), "Unexpected valid chain")
}

func TestAlteringIndexInvalidatesChain(t *testing.T) {
   blockchain := New()
   lastProof := blockchain.LastBlock().Proof
   proof := ProofOfWork(lastProof)
   blockchain.NewBlock(proof)

   // Alter the previous index block
   blockchain.Chain()[0].Index = 666

   assert.NotNil(t, blockchain.Validate(), "Unexpected valid chain")
}

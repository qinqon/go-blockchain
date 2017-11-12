package blockchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	assert.Error(t, blockchain.Validate(), "Unexpected valid chain")
}

func TestAlteringIndexInvalidatesChain(t *testing.T) {
	blockchain := New()
	lastProof := blockchain.LastBlock().Proof
	proof := ProofOfWork(lastProof)
	blockchain.NewBlock(proof)

	// Alter the previous index block
	blockchain.Chain()[0].Index = 666

	assert.Error(t, blockchain.Validate(), "Unexpected valid chain")
}

func TestAlteringTransactionInvalidatesChain(t *testing.T) {
	blockchain := New()
	proof := blockchain.LastBlock().Proof

	// Add a transaction
	blockchain.NewTransaction(Transaction{Sender: "qinqon", Recipient: "golang", Amount: 99999})

	proof = ProofOfWork(proof)
	blockchain.NewBlock(proof)

	proof = ProofOfWork(proof)
	blockchain.NewBlock(proof)

	// Alter the transaction
	blockchain.Chain()[1].Transactions[0].Recipient = "Bad guy"

	assert.Error(t, blockchain.Validate(), "Unexpected valid chain")
}

func TestNoConflictsWithSmallerBlockchain(t *testing.T) {
	blockchain := New()
	proof := blockchain.LastBlock().Proof

	proof = ProofOfWork(proof)
	blockchain.NewBlock(proof)

	proof = ProofOfWork(proof)
	blockchain.NewBlock(proof)

	conflictingBlockchain := New()
	proof = conflictingBlockchain.LastBlock().Proof

	proof = ProofOfWork(proof)
	conflictingBlockchain.NewBlock(proof)

	expectedChain := blockchain.Chain()

	blockchain.ResolveConflict(conflictingBlockchain.Chain())
	assert.Equal(t, expectedChain, blockchain.Chain())
}

func TestReplaceConflictResolutionWithBiggerBlockchain(t *testing.T) {
	blockchain := New()
	proof := blockchain.LastBlock().Proof

	proof = ProofOfWork(proof)
	blockchain.NewBlock(proof)

	conflictingBlockchain := New()
	proof = conflictingBlockchain.LastBlock().Proof

	proof = ProofOfWork(proof)
	conflictingBlockchain.NewBlock(proof)

	proof = ProofOfWork(proof)
	conflictingBlockchain.NewBlock(proof)

	expectedChain := conflictingBlockchain.Chain()

	blockchain.ResolveConflict(conflictingBlockchain.Chain())
	assert.Equal(t, expectedChain, blockchain.Chain())
}

func TestNoConflictWithInvalidBlockchain(t *testing.T) {
	blockchain := New()
	proof := blockchain.LastBlock().Proof
	proof = ProofOfWork(proof)
	blockchain.NewBlock(proof)

	conflictingBlockchain := New()
	proof = conflictingBlockchain.LastBlock().Proof
	proof = ProofOfWork(proof)
	conflictingBlockchain.NewBlock(proof)
	proof = ProofOfWork(proof)
	conflictingBlockchain.NewBlock(proof)

	// Alter the previous index block
	conflictingBlockchain.Chain()[0].Index = 666

	expectedChain := blockchain.Chain()

	blockchain.ResolveConflict(conflictingBlockchain.Chain())
	assert.Equal(t, expectedChain, blockchain.Chain())
}

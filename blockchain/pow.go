package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func ValidProof(lastProof, proof int) error {
	guess := fmt.Sprintf("%v%v", lastProof, proof)
	guessHash := sha256.Sum256([]byte(guess))
	guessHashEncoded := hex.EncodeToString(guessHash[:])
	if !(guessHashEncoded[len(guessHashEncoded)-4:] == "0000") {
		return fmt.Errorf("Hash '%v' of '%v' not ending with 4 zeroes", guessHashEncoded, guess)
	}
	return nil
}

func ProofOfWork(lastProof int) int {
	proof := 0
	for ValidProof(lastProof, proof) != nil {
		proof++
	}
	return proof
}

package blockchain

import(
	"errors"
	"time"
	"encoding/json"
	"crypto/sha256"
	"encoding/hex"
)


type Block struct {
   Index int `json:"index"`
   Timestamp time.Time `json:"timestamp"`
   Transactions []Transaction `json:"transactions"`
   Proof int `json:"proof"`
   PreviousHash string `json:"previous_hash"`
}

func (block *Block) CalculateHash() (string, error){
   //TODO: Sort dictionary
   marshalledBlock, err := json.Marshal(*block)
   if err != nil {
      return "", err
   }
   hashedBlock := sha256.Sum256(marshalledBlock)
   hexBlock := hex.EncodeToString(hashedBlock[:])
   return hexBlock, nil
}

func (block *Block) Validate(previousBlock *Block) error {
	if err := ValidProof(previousBlock.Proof, block.Proof); err != nil{
		return err;
	}

	previousBlockHash, err := previousBlock.CalculateHash()
	if err != nil {
		return err
	}
	if block.PreviousHash != previousBlockHash {
		return errors.New("Invalid previous hash")
	}
	return nil
}

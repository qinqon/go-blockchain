package blockchain

import "time"

type Block struct {
   Index int `json:"index"`
   Timestamp time.Time `json:"timestamp"`
   Transactions []Transaction `json:"transactions"`
   Proof int `json:"proof"`
   PreviousHash string `json:"previous_hash"`
}

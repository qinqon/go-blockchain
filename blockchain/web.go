package blockchain

import (
   "fmt"
   "net/http"
   "encoding/json"
   "io/ioutil"
   "github.com/satori/go.uuid"
)

var node_identifier = uuid.NewV4().String()


func (b *Blockchain) mineHandler(w http.ResponseWriter, r *http.Request) {
   last_proof := b.LastBlock().Proof
   proof := ProofOfWork(last_proof)
   b.NewTransaction(Transaction{"0", node_identifier, 1})
   block := b.NewBlock(proof)
   w.Header().Add("Content-Type", "application/json")
   fmt.Fprintf(w, `{
      "message":"New Block Forged", 
      "index": %v,
      "tarnsactions": %v,
      "proof": %v,
      "previous_hash": %v,
   }`, block.Index, block.Transactions, block.Proof, block.PreviousHash)

}


func (b *Blockchain) newTransactionHandler(w http.ResponseWriter, r *http.Request) {
   defer r.Body.Close()
   body, err := ioutil.ReadAll(r.Body)
   if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
   }
   transaction := Transaction{}
   err = json.Unmarshal(body, &transaction)
   if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
   }
   index := b.NewTransaction(transaction)
   w.Header().Add("Content-Type", "application/json")
   fmt.Fprintf(w, `{"message":"Transaction will be added to Block %v"}`, index)
}

func (b *Blockchain) fullChainHandler(w http.ResponseWriter, r *http.Request) {
   marshalledChain, err := json.Marshal(b.Chain())
   if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
   }
   w.Header().Add("Content-Type", "application/json")
   fmt.Fprintf(w, string(marshalledChain))
}

func (b *Blockchain) StartWebServer() {
   http.HandleFunc("/mine", b.mineHandler)
   http.HandleFunc("/transaction/new", b.newTransactionHandler)
   http.HandleFunc("/chain", b.fullChainHandler)
   http.ListenAndServe(":5000", nil)
}

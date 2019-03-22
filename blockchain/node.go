package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Node struct {
	neighbours map[*url.URL]bool
	blockchain Blockchain
	identifier string
	httpServer *http.Server
}

func NewNode(address string) *Node {
	u, err := uuid.NewV4()
	if err != nil {
		return nil
	}
	node := Node{
		identifier: u.String(),
	}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/mine", node.mineHandler)
	serveMux.HandleFunc("/transaction/new", node.newTransactionHandler)
	serveMux.HandleFunc("/chain", node.fullChainHandler)
	node.httpServer = &http.Server{
		Addr:    address,
		Handler: serveMux,
	}
	return &node
}

func (n *Node) RegisterNeighbour(address string) error {
	u, err := url.Parse(address)
	if err != nil {
		return err
	}
	n.neighbours[u] = true
	return nil
}

func (n *Node) ResolveConflicts() {
	for neighbour := range n.neighbours {
		resp, err := http.Get("http://" + neighbour.String() + "/chain")
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		var neighbourChain []*Block
		if err = json.NewDecoder(resp.Body).Decode(neighbourChain); err != nil {
			continue
		}
		n.blockchain.ResolveConflict(neighbourChain)
	}
}

func (n *Node) mineHandler(w http.ResponseWriter, r *http.Request) {
	last_proof := n.blockchain.LastBlock().Proof
	proof := ProofOfWork(last_proof)
	n.blockchain.NewTransaction(Transaction{"0", n.identifier, 1})
	block := n.blockchain.NewBlock(proof)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{
      "message":"New Block Forged", 
      "index": %v,
      "tarnsactions": %v,
      "proof": %v,
      "previous_hash": %v,
   }`, block.Index, block.Transactions, block.Proof, block.PreviousHash)

}

func (n *Node) newTransactionHandler(w http.ResponseWriter, r *http.Request) {
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
	index := n.blockchain.NewTransaction(transaction)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message":"Transaction will be added to Block %v"}`, index)
}

func (n *Node) fullChainHandler(w http.ResponseWriter, r *http.Request) {
	marshalledChain, err := json.Marshal(n.blockchain.Chain())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(marshalledChain))
}

func (n *Node) HttpServer() *http.Server {
	return n.httpServer
}

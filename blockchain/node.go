package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Node struct {
	neighbours map[*url.URL]bool
	blockchain *Blockchain
	identifier string
	httpServer *http.Server
}

func NewNode(address string) *Node {
	node := Node{
		identifier: uuid.NewV4().String(),
	}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/mine", node.mineHandler)
	serveMux.HandleFunc("/transaction/new", node.newTransactionHandler)
	serveMux.HandleFunc("/chain", node.fullChainHandler)
	serveMux.HandleFunc("/neighbour", node.registerNeighbourHandler)
	serveMux.HandleFunc("/resolve", node.resolveConflictsHandler)
	node.httpServer = &http.Server{
		Addr:    address,
		Handler: serveMux,
	}
	node.blockchain = NewBlockchain()
	node.neighbours = make(map[*url.URL]bool)
	return &node
}

func (node *Node) Start(address string) error {
	log.Printf("Starting blockchain node at %v", address)
	if err := node.HttpServer().ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (n *Node) RegisterNeighbour(address string) error {
	u, err := url.Parse(address)
	if err != nil {
		return err
	}
	n.neighbours[u] = true
	return nil
}

func (n *Node) ResolveConflicts() bool {
	chainChanged := false
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
		if n.blockchain.ResolveConflict(neighbourChain) {
			chainChanged = true
		}
	}
	return chainChanged
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

func (n *Node) registerNeighbourHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nodes := Nodes{}
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, node := range nodes.Nodes {
		n.RegisterNeighbour(node)
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message":"neighbours added"}`)
}

func (n *Node) resolveConflictsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	chainReplaced := n.ResolveConflicts()
	w.Header().Add("Content-Type", "application/json")
	if chainReplaced {
		fmt.Fprintf(w, `{"message":"chain replaced"}`)
	} else {
		fmt.Fprintf(w, `{"message":"chain maintened"}`)
	}
}
func (n *Node) HttpServer() *http.Server {
	return n.httpServer
}

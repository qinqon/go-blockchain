package blockchain

import (
   "fmt"
   "net/http"
)

func mineHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Mine")
}


func newTransactionHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "New transaction")
}

func fullChainHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Full chain")
}


func StartWebServer() {
   http.HandleFunc("/mine", mineHandler)
   http.HandleFunc("/transaction/new", newTransactionHandler)
   http.HandleFunc("/chain", fullChainHandler)
   http.ListenAndServe(":8080", nil)
}

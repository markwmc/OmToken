package main

import (
	"fmt"
	"log"
	"net/http"


	"github.com/gorilla/mux"
	"github.com/markwmc/OmToken/blockchain"
)

func main() {
	blockchain.InitBlockchain()
	blockchain.InitializeWallets()

	router := mux.NewRouter()

	 router.HandleFunc("/addTransaction", blockchain.AddTransaction).Methods("POST")

	fmt.Println("Blockchain API running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
	
}
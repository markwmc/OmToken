package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	

	blockchain "github.com/markwmc/OmToken/blockchain"
)

var Wallets []*blockchain.Wallet
var Blockchain []*blockchain.Block
var TransactionPool []blockchain.Transaction


func CreateWallet(w http.ResponseWriter, r *http.Request) {
	newWallet := blockchain.GenerateWallet("test")
	Wallets = append(Wallets, newWallet)

	response := map[string]string{
		"Address":    newWallet.Address,
		"PublicKey":  fmt.Sprintf("%x", newWallet.PublicKey),
		"PrivateKey": fmt.Sprintf("%x", newWallet.PrivateKey.D.Bytes()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SignTransaction(w http.ResponseWriter, r *http.Request) {
	if len(Wallets) == 0 {
		http.Error(w, "No wallets available", http.StatusBadRequest)
		return
	}

	message := "Alice pays Bob 2 OT"
	wallet := Wallets[0]

	rSig, sSig := wallet.SignTransaction(message)

	response := map[string]string{
		"Message": message,
		"R":       rSig,
		"S":       sSig,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func VerifyTransaction(w http.ResponseWriter, r *http.Request) {
	if len(Wallets) == 0 {
		http.Error(w, "No wallets available", http.StatusBadRequest)
		return
	}

	message := "Alice pays Bob 2 OT"
	wallet := Wallets[0]

	rSig, sSig := wallet.SignTransaction(message)

	pubKeyBytes := append(wallet.PublicKey.X.Bytes(), wallet.PublicKey.Y.Bytes()...)

	valid := blockchain.VerifySignature(pubKeyBytes, message, rSig, sSig)

	response := map[string]bool{"validSignature": valid}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func MineBlock(w http.ResponseWriter, _ *http.Request) {
	transactions := []string{}
	for _, tx := range TransactionPool {
		transactions = append(transactions, fmt.Sprintf("%s -> %s: %.2f OT", tx.Sender, tx.Receiver, tx.Amount))
	}
	if len(blockchain.Blockchain) == 0 {
		http.Error(w, "Blockchain not initialized", http.StatusInternalServerError)
		return
	}


	prevBlock := blockchain.Blockchain[len(blockchain.Blockchain)-1]
	newBlock := blockchain.GenerateBlock(prevBlock, transactions)
	blockchain.Blockchain = append(blockchain.Blockchain, newBlock)

	TransactionPool = nil

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(newBlock); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Block mined successfully: %s\n", newBlock.Hash)
}

func GetBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockchain.Blockchain)
}

func AddBlock(w http.ResponseWriter, r *http.Request) {
	transactions := []string{"Alice -> Bob: 2 OT", "Bob -> Charlie: 1 OT"}
	blockchain.AddBlock(transactions)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Block added"})
}



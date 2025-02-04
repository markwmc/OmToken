package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sort"
)

type Transaction struct {
	Sender		string
	Receiver	string
	Amount		float64
	Fee			float64
	Signature	string
}

var userBalances = map[string]float64{
	"Alice": 10.0,
	"Bob": 5.0,
	"Charlie": 2.0,
	"Dave": 0.5,
}

func validateTransaction(tx Transaction, senderPublicKey *ecdsa.PublicKey) bool {
	senderBalance, exists := userBalances[tx.Sender]

	if !exists {
		fmt.Println("Sender not found")
		return false
	}
	if senderBalance < tx.Amount+tx.Fee {
		fmt.Println("Insufficient funds for the transaction")
		return false
	}

	if !verifySignature(tx, senderPublicKey){
		fmt.Println("Invalid signature")
		return false
	}
	if tx.Fee <= 0 {
		fmt.Println("Transaction fee must be greater than zero")
		return false
	}

	if tx.Receiver == "" {
		fmt.Println("Invalid receiver address")
		return false
	}

	return true

}

func verifySignature(tx Transaction, pubKey *ecdsa.PublicKey) bool {
	txData := fmt.Sprintf("%s%s%f%f", tx.Sender, tx.Receiver, tx.Amount, tx.Fee)
	hash := sha256.Sum256([]byte(txData))

	signatureBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		fmt.Println("Error decoding signature:", err)
		return false
	}

	r := new(big.Int).SetBytes(signatureBytes[:len(signatureBytes)/2])
	s := new(big.Int).SetBytes(signatureBytes[len(signatureBytes)/2:])

	return ecdsa.Verify(pubKey, hash[:], r, s)
}

func getPublicKeyFromWallet(sender string) *ecdsa.PublicKey {
	for _, wallet := range Wallets {
		if wallet.Address == sender {
			return wallet.PublicKey
		}
	}
	return nil
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionPool []Transaction
	var tx Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	publicKey := getPublicKeyFromWallet(tx.Sender)
	if publicKey == nil {
		http.Error(w, "Sender not found", http.StatusBadRequest)
		return
	}

	if !validateTransaction(tx, publicKey) {
		http.Error(w, "Invalid transaction", http.StatusBadRequest)
		return
	}

	transactionPool = append(transactionPool, tx)
	sort.Slice(transactionPool, func(i, j int) bool {
		return transactionPool[i].Fee > transactionPool[j].Fee
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction added"})
}
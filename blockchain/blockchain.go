package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	
)


const Difficulty = 4


var Blockchain []*Block


func GenerateBlock(prevBlock *Block, transactions []string) *Block {
	newBlock := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		PrevHash:     prevBlock.Hash,
		Transactions: transactions,
	}
	newBlock.MerkleRoot = BuildMerkleTree(newBlock.Transactions)
	newBlock.Hash = MineBlock(newBlock)

	return newBlock
}

func InitBlockchain() {
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		PrevHash:     "",
		Transactions: []string{"Genesis Block"},
	}
	genesisBlock.MerkleRoot = BuildMerkleTree(genesisBlock.Transactions)
	genesisBlock.Hash = MineBlock(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)
}


func GenerateHash(block *Block) string {
	record := fmt.Sprintf("%d%s%s%s%d", block.Index, block.Timestamp, block.PrevHash, block.MerkleRoot, block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}


func MineBlock(block *Block) string {
	for {
		hash := GenerateHash(block)
		if hash[:Difficulty] == string(bytes.Repeat([]byte{'0'}, Difficulty)) {
			block.Hash = hash
			fmt.Printf("Block Mined! Nonce: %d, Hash: %s\n", block.Nonce, block.Hash)
			return hash
		}
		block.Nonce++
	}
}

func AddBlock(transactions []string) {
	prevBlock := Blockchain[len(Blockchain)-1]
	newBlock := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		PrevHash:     prevBlock.Hash,
		Transactions: transactions,
	}
	newBlock.MerkleRoot = BuildMerkleTree(newBlock.Transactions)
	newBlock.Hash = MineBlock(newBlock)
	Blockchain = append(Blockchain, newBlock)
	fmt.Println("Block successfully added")
}

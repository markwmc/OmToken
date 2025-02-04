package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
)


type Block struct {
	Index        int
	Timestamp    string
	PrevHash     string
	Hash         string
	Transactions []string
	MerkleRoot   string
	Nonce        int
}


type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  string
}


func NewMerkleNode(left, right *MerkleNode, data string) *MerkleNode {
	var hash string
	if left == nil && right == nil {
		hash = hashData(data)
	} else {
		combined := left.Hash + right.Hash
		hash = hashData(combined)
	}
	return &MerkleNode{Left: left, Right: right, Hash: hash}
}


func hashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}


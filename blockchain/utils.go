package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
)

func BuildMerkleTree(transactions []string) string {
	if len(transactions) == 0 {
		return ""
	}

	var nodes []*MerkleNode
	for _, tx := range transactions {
		nodes = append(nodes, NewMerkleNode(nil, nil, tx))
	}

	for len(nodes) > 1 {
		var newLevel []*MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			if i+1 == len(nodes) {
				nodes = append(nodes, nodes[i])
			}
			newLevel = append(newLevel, NewMerkleNode(nodes[i], nodes[i+1], ""))
		}
		nodes = newLevel
	}

	return nodes[0].Hash
}

func sha256Hash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
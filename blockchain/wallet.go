package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
}

var Wallets = []Wallet{}


func InitializeWallets() {
	GenerateWallet("Alice")
	GenerateWallet("Bob")
	GenerateWallet("Charlie")
	GenerateWallet("Dave")

}
func GenerateWallet(name string) *Wallet {
	privKey, pubKey, err := generateKeys()
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return nil
	}

	address := GenerateAddress(pubKey)

	return &Wallet{
		PrivateKey: privKey,
		PublicKey: pubKey,
		Address: address,
	}
}

func GenerateAddress(pubKey *ecdsa.PublicKey) string {
	pubBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	hash := sha256.Sum256(pubBytes)
	address := hex.EncodeToString(hash[:])
	return address
}

func (w *Wallet) SignTransaction(message string) (string, string) {
	hash := sha256.Sum256([]byte(message))
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		fmt.Println("Error signing transaction:", err)
		return "", ""
	}

	return r.Text(16), s.Text(16)
}

func VerifySignature(publicKey []byte, message string, rStr, sStr string) bool {
	hash := sha256.Sum256([]byte(message))

	curve := elliptic.P256()
	x, y := big.Int{}, big.Int{}
	keyLen := len(publicKey) / 2
	x.SetBytes(publicKey[:keyLen])
	y.SetBytes(publicKey[keyLen:])

	r, _ := new(big.Int).SetString(rStr, 16)
	s, _ := new(big.Int).SetString(sStr, 16)

	pubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	return ecdsa.Verify(&pubKey, hash[:], r, s)
}

func generateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

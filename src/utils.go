package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

// RSAbitSize ... Bit size for key pair
const (
	RSAbitSize      = 10
	MiningDificulty = 2
	MiningPrize     = 1
)

func generateSHA256Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}

func generateKeyPair() *rsa.PrivateKey {
	reader := rand.Reader
	key, _ := rsa.GenerateKey(reader, RSAbitSize)
	return key
}

func createTransaction(senderID string, receiverID string, amount float32) (transaction, error) {
	debug(fmt.Sprintf("Creating transaction - senderID:%v receiverID:%v amount:%v", senderID, receiverID, amount))

	if amount <= 0 {
		return transaction{}, fmt.Errorf("Wrong Amount Passed %v", amount)
	}

	if receiverID == "" {
		return transaction{}, fmt.Errorf("Empty receiver ID")
	}

	if senderID == receiverID {
		return transaction{}, fmt.Errorf("Sender and receiver id cannot be same %v", senderID)
	}

	return transaction{senderID: senderID, receiverID: receiverID, amount: amount, timestamp: time.Now()}, nil
}

func debug(s string) {
	if debugP {
		fmt.Println(s)
	}
}

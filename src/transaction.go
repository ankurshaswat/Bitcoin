package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"
)

type transaction struct {
	senderID, receiverID string
	amount               float32
	timestamp            time.Time
	signature            []byte
}

func (t *transaction) getHash() string {
	amountString := fmt.Sprintf("%f", t.amount)
	s := t.senderID + t.receiverID + amountString
	hash := generateSHA256Hash(s)
	return hash
}

func (t *transaction) signTransaction(keyPair *rsa.PrivateKey) {
	// TODO: check if public key is senders public key (probab write a function for node to get its public key)

	// Get hash of transaction and sign it with key
	hash := t.getHash()

	rng := rand.Reader
	signature, err := rsa.SignPKCS1v15(rng, keyPair, crypto.SHA256, []byte(hash))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return
	}

	// Return after saving signature in transaction
	t.signature = signature
	return
}

func (t *transaction) verifyTransaction() bool {

	//  No from address means mining reward self added
	// ? Are any further checks required here
	if t.senderID == "" {
		return true
	}

	if t.signature == nil {
		log.Fatal("No signature found")
	}

	// TODO: get public key of senderID
	// pubKey :=

	// TODO: verify signature with hash of transaction
	// err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, t.getHash(), t.signature)
	// if err != nil {
	// return false
	// }

	return true
}

func createTransaction(senderID string, receiverID string, amount float32) transaction {
	return transaction{senderID: senderID, receiverID: receiverID, amount: amount, timestamp: time.Now()}
}
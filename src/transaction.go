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
	// check if public key is senders public key (probab write a function for node to get its public key)
	senderID := t.senderID
	pubKeySender := getPublicKey(senderID)
	if *pubKeySender != keyPair.PublicKey {
		log.Fatal("Incosistent public key found while signing")
	}

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

func (t *transaction) verifyTransaction() (bool, error) {

	//  No from address means mining reward self added
	// ? Are any further checks required here
	if t.senderID == "" {
		return false, fmt.Errorf("Sender Id is empty")
	}

	if t.signature == nil {
		return false, fmt.Errorf("No signature found")
	}

	// get public key of senderID
	pubKey := getPublicKey(t.senderID)

	// verify signature with hash of transaction
	err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, []byte(t.getHash()), t.signature)
	if err != nil {
		return false, err
	}

	return true, nil
}

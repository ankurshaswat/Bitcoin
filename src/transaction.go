package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"
)

type transaction struct {
	senderID, receiverID string
	amount               float64
	timestamp            time.Time
	signature            []byte
}

func (t *transaction) getHash() string {
	amountString := fmt.Sprintf("%f", t.amount)
	s := t.senderID + t.receiverID + amountString + strconv.FormatInt(t.timestamp.UnixNano(), 10)
	hash := generateSHA256Hash(s)
	return hash
}

func (t *transaction) signTransaction(keyPair *rsa.PrivateKey) {
	// check if public key is senders public key (probab write a function for node to get its public key)
	senderID := t.senderID
	pubKeySender := getPublicKey(senderID)
	if *pubKeySender != keyPair.PublicKey {
		log.Panic("Incosistent public key found while signing")
	}

	// Get hash of transaction and sign it with key
	hash := t.getHash()
	hashInBytes, err := hex.DecodeString(hash)
	if err != nil {
		log.Panicf("Error in decoding hex %v err:%v", hash, err)
	}
	// log.Println(hash, hashInBytes)
	rng := rand.Reader
	signature, err := rsa.SignPKCS1v15(rng, keyPair, crypto.SHA256, hashInBytes)
	if err != nil {
		log.Panicf("Error from signing: %s\n", err)
	}

	// Return after saving signature in transaction
	t.signature = signature
	return
}

func (t *transaction) verifyTransaction() (bool, error) {

	//  No from address means mining reward self added
	// ? Are any further checks required here
	if t.senderID == "" {
		return true, nil
		// return false, fmt.Errorf("Sender Id is empty")
	}

	if t.signature == nil {
		return false, fmt.Errorf("No signature found")
	}

	// get public key of senderID
	pubKey := getPublicKey(t.senderID)

	hash := t.getHash()
	hashInBytes, err := hex.DecodeString(hash)
	if err != nil {
		log.Panicf("Error in decoding hex %v err:%v", hash, err)
	}
	// verify signature with hash of transaction
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashInBytes, t.signature)
	if err != nil {
		return false, err
	}

	return true, nil
}

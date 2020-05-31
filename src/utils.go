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

func createMerkleTree(tList []transaction) treeBlock {

	blockList := []treeBlock{}

	for i := 0; i < len(tList); i += 2 {
		trans1 := tList[i]
		if i+1 == len(tList) {
			// If only one left
			hash := generateSHA256Hash(trans1.getHash())
			newLeafBlock := treeBlock{leaf: true, leftT: &trans1, hash: hash}
			blockList = append(blockList, newLeafBlock)
		} else {
			// if more than one available
			trans2 := tList[i+1]
			hash := generateSHA256Hash(trans1.getHash() + trans2.getHash())
			newLeafBlock := treeBlock{leaf: true, leftT: &trans1, rightT: &trans2, hash: hash}
			blockList = append(blockList, newLeafBlock)
		}
	}

	for len(blockList) > 1 {
		newBlockList := []treeBlock{}

		for i := 0; i < len(blockList); i += 2 {
			block1 := blockList[i]
			if i+1 == len(blockList) {
				// If only one left
				hash := generateSHA256Hash(block1.hash)
				newTreeBlock := treeBlock{leaf: false, left: &block1, hash: hash}
				newBlockList = append(newBlockList, newTreeBlock)
			} else {
				// if two available
				block2 := blockList[i+1]
				hash := generateSHA256Hash(block1.hash + block2.hash)
				newTreeBlock := treeBlock{leaf: false, left: &block1, right: &block2, hash: hash}
				newBlockList = append(newBlockList, newTreeBlock)
			}
		}
		blockList = newBlockList
	}

	return blockList[0]
}

func createBlock(transactions []transaction, prevHash string) block {
	transactionTree := createMerkleTree(transactions)
	return block{timestamp: time.Now(), prevHash: prevHash, transactionTree: transactionTree}
}

func createNode(nodeID string, receiveChan chan block, blockchain []block) node {
	keyPair := generateKeyPair()
	n := node{nodeID: nodeID, keyPair: keyPair, receiveChannel: receiveChan, blockchain: blockchain}
	return n
}

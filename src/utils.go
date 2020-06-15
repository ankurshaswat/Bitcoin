package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

func generateSHA256Hash(s string) string {
	message := []byte(s)
	hashInBytes := sha256.Sum256(message)

	// h := sha1.New()
	// h.Write([]byte(s))
	// hashInBytes := h.Sum(nil)
	// log.Println(hashInBytes)
	// log.Println(len(hashInBytes))
	sha1Hash := hex.EncodeToString(hashInBytes[:])
	// log.Println(sha1Hash)
	return sha1Hash
}

func generateKeyPair() *rsa.PrivateKey {
	reader := rand.Reader
	key, err := rsa.GenerateKey(reader, rsaBitSize)
	if err != nil {
		log.Panic("Unable to generate Rsa key - ", err)
	}
	return key
}

func debug(s string) {
	if debugP {
		log.Println(s)
	}
}

func createMerkleTree(tList []transaction) merkleTree {

	blockList := []merkleTree{}

	for i := 0; i < len(tList); i += 2 {
		trans1 := tList[i]
		if i+1 == len(tList) {
			// If only one left
			hash := generateSHA256Hash(trans1.hash)
			newLeafBlock := merkleTree{leaf: true, leftT: &trans1, hash: hash}
			blockList = append(blockList, newLeafBlock)
		} else {
			// if more than one available
			trans2 := tList[i+1]
			hash := generateSHA256Hash(trans1.hash + trans2.hash)
			newLeafBlock := merkleTree{leaf: true, leftT: &trans1, rightT: &trans2, hash: hash}
			blockList = append(blockList, newLeafBlock)
		}
	}

	for len(blockList) > 1 {
		newBlockList := []merkleTree{}

		for i := 0; i < len(blockList); i += 2 {
			block1 := blockList[i]
			if i+1 == len(blockList) {
				// If only one left
				hash := generateSHA256Hash(block1.hash)
				newTreeBlock := merkleTree{leaf: false, left: &block1, hash: hash}
				newBlockList = append(newBlockList, newTreeBlock)
			} else {
				// if two available
				block2 := blockList[i+1]
				hash := generateSHA256Hash(block1.hash + block2.hash)
				newTreeBlock := merkleTree{leaf: false, left: &block1, right: &block2, hash: hash}
				newBlockList = append(newBlockList, newTreeBlock)
			}
		}
		blockList = newBlockList
	}

	return blockList[0]
}

func extractBalTransaction(tx *transaction, nodeID string) float64 {
	bal := 0.0
	if tx.senderID == nodeID {
		bal -= tx.amount
	} else if tx.receiverID == nodeID {
		bal += tx.amount
	}
	return bal
}

func calcBalance(tree *merkleTree, nodeID string) float64 {
	sum := 0.0

	if tree.leaf {
		sum += extractBalTransaction(tree.leftT, nodeID)
		if tree.rightT != nil {
			sum += extractBalTransaction(tree.rightT, nodeID)
		}
	} else {
		sum += calcBalance(tree.left, nodeID)
		if tree.right != nil {
			sum += calcBalance(tree.right, nodeID)
		}
	}

	return sum
}

func createBalanceCopy(original map[string]float64) map[string]float64 {
	newMap := make(map[string]float64)

	// Copy from the original map to the target map
	for key, value := range original {
		newMap[key] = value
	}

	return newMap
}

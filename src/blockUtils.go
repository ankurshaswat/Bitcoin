package main

import (
	"log"
	"time"
)

func createGenesisBlock() block {
	var transactionTree merkleTree
	b := block{timestamp: time.Now(), prevHash: "", transactionTree: transactionTree}
	b.mine()
	return b
}

func createBlock(transactions []transaction, prevHash string) block {
	transactionTree := createMerkleTree(transactions)
	b := block{timestamp: time.Now(), prevHash: prevHash, transactionTree: transactionTree}
	return b
}

func broadcastBlock(b block, selfID string) {
	// TODO: Maybe randomize this to send to less number of nodes
	log.Println("Node:", selfID, "queing request for new block.")

	for i := 0; i < len(nodeList); i++ {
		nodeInstance := &nodeList[i]
		if nodeInstance.nodeID != selfID {
			// Append Block in new goRoutine
			go nodeInstance.addBlock(b, selfID)
			// nodeInstance.receiveChannel <- b
		}
	}
	log.Println("Node:", selfID, "published new block. PrevHash:", b.prevHash[:15], "NewHash:", b.hash[:15])
}

package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"strconv"
)

func getPublicKey(nodeID string) *rsa.PublicKey {
	i, err := strconv.ParseInt(nodeID, 10, 64)
	if err != nil {
		fmt.Print(err)
		log.Fatal("Error in Prasing int")
	}
	return nodeList[i].getPublicKey()
}

func broadcastBlock(b block, selfID string) {
	// TODO: Maybe randomize this to send to less number of nodes

	for i := 0; i < len(nodeList); i++ {
		nodeInstance := nodeList[i]
		if nodeInstance.nodeID != selfID {
			nodeInstance.receiveChannel <- b
		}
	}
}

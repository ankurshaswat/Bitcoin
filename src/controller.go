package main

import (
	"crypto/rsa"
	"log"
	"strconv"
)

func getPublicKey(nodeID string) *rsa.PublicKey {
	i, err := strconv.ParseInt(nodeID, 10, 64)
	if err != nil {
		log.Print(err)
		log.Panic("Error in Prasing int")
	}
	return nodeList[i].getPublicKey()
}

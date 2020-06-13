package main

func createNode(nodeID string, blockchain []block, nodeType int) node {
	keyPair := generateKeyPair()
	// n :=
	return node{nodeID: nodeID, keyPair: keyPair, blockchain: blockchain, nodeType: nodeType}
}

package main

import (
	"fmt"
	"strconv"
)

var debugP = false
var nodeList []node
var recChanList []chan block

func main() {

	// Writing Tests
	// runTests()

	// Create genesis block
	genesisBlock := createBlock([]transaction{}, "")
	fmt.Println("Genesis block ", genesisBlock)

	numNodes := 10

	// Create new nodes in loop (while also providing genesis block ref)
	for i := 0; i < numNodes; i++ {
		recChan := make(chan block)
		recChanList = append(recChanList, recChan)

		// ? Do we always pass only genesis block to all nodes in network
		n := createNode(strconv.Itoa(i), recChan, []block{genesisBlock}, normal)
		nodeList = append(nodeList, n)
		go n.startNode()
		fmt.Printf("Started node number %v \n", i+1)
	}

	// TODO: Order nodes to do transactions
}

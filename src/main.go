package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

var debugP = false
var nodeList []node
var startMining = false

const (
	rsaBitSize      = 2048
	miningDificulty = 5
	miningPrize     = 1
	numNodes        = 10
)

// var recChanList []chan block

func main() {

	log.SetFlags(log.Lmicroseconds)
	// Writing Tests
	// runTests()

	// Create genesis block
	genesisBlock := createGenesisBlock()
	log.Println("Genesis block hash:", genesisBlock.hash[:15])

	// Create new nodes in loop (while also providing genesis block ref)
	for i := 0; i < numNodes; i++ {
		// recChan := make(chan block)
		// recChanList = append(recChanList, recChan)

		// ? Do we always pass only genesis block to all nodes in network
		// n :=
		nodeList = append(nodeList, createNode(strconv.Itoa(i), []block{genesisBlock}, normal))
	}

	for i := 0; i < numNodes; i++ {
		log.Printf("Starting node number %v \n", i+1)
		go nodeList[i].startNode()
	}

	startMining = true

	for true {
		// log.Print("Insert Command here: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		log.Println(input.Text())
	}

	// TODO: Order nodes to do transactions
}

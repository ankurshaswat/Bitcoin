package main

import (
	"crypto/rsa"
	"log"
	"sync"
	"time"
)

// Node type
const (
	normal           = iota
	malicious        = iota
	smartSingleUse   = iota
	smartRepeatedUse = iota
)

type node struct {
	sync.Mutex
	nodeID     string
	blockchain []block
	// blockchainLock      sync.Mutex
	pendingTransactions []transaction // to be pushed onto next block
	// pendingTxLock       sync.Mutex
	keyPair        *rsa.PrivateKey
	smartBalTarget float64
	nodeType       int
	balances       map[string]float64
}

// startMiner starts an endless loop which keeps on searching for a new nonce and keeps broadcasting new blocks.
func (n *node) startMiner() {

	for true {
		n.Lock()
		// n.pendingTxLock.Lock()
		// n.blockchainLock.Lock()

		// Create a new transaction giving yourself reward and add to pending transactions
		tx, err := createTransaction("", n.nodeID, miningPrize)
		if err != nil {
			log.Panic("Error mining coin - creating transaction -", err)
		}
		n.pendingTransactions = append(n.pendingTransactions, tx)

		// Create block with all pending transactions and mine to get Nonce
		// log.Println(n.nodeID, n.blockchain)
		// if len(n.blockchain) > 1 {

		// 	log.Println(len(n.blockchain))
		// }
		prevBlockHash := n.blockchain[len(n.blockchain)-1].hash
		b := createBlock(n.pendingTransactions, prevBlockHash)
		mineSuccess := b.mineSingleTry()

		if !mineSuccess {
			// Remove last transaction from pendingTransactions
			n.pendingTransactions = n.pendingTransactions[:len(n.pendingTransactions)-1]
		} else {
			for _, tx := range n.pendingTransactions {
				n.balances[tx.receiverID] += tx.amount
				if tx.senderID != "" {
					n.balances[tx.senderID] -= tx.amount
				}
				// n.selfBal += extractBalTransaction(&tx, n.nodeID)
			}

			// Push block onto chain
			n.blockchain = append(n.blockchain, b)
			n.pendingTransactions = []transaction{}

			// log.Println(len(n.blockchain))
			// Broadcast new block to all nodes
			broadcastBlock(b, n.nodeID)
			// log.Println("Node:", n.nodeID, "last hash", n.blockchain[len(n.blockchain)-1].hash, len(n.blockchain))
		}

		// n.pendingTxLock.Unlock()
		// n.blockchainLock.Unlock()
		n.Unlock()
	}
}

func (n *node) addTransaction(tx transaction) {
	// Add transaction to pending Transactions list after checking validity and possibility and existence etc.
	if tx.senderID == "" || tx.receiverID == "" {
		log.Panic("Empty sender or receiver ID")
	}

	txverified, err := tx.verifyTransaction()
	if err != nil {
		log.Panic("Error in verifiying transaction - ", err)
	}
	if !txverified {
		log.Panic("Transaction verification failed")
	}

	if tx.amount <= 0 {
		log.Panic("Amount <= 0")
	}

	if n.balances[tx.senderID] < tx.amount {
		log.Panic("Insufficient balance for transaction")
	}

	n.pendingTransactions = append(n.pendingTransactions, tx)
}

// func (n *node) getBalance(nodeID string) float64 {
// 	if nodeID == n.nodeID {
// 		return n.selfBal
// 	}

// 	bal := 0.0
// 	// Go over all blocks and all transactions to generate final balance of a node
// 	for _, b := range n.blockchain {
// 		bal += calcBalance(&b.transactionTree, nodeID)
// 	}

// 	return bal
// }

func (n *node) verifyChain() bool {
	// ? Add check of genesis block
	// TODO: Add checks of balances here

	// Verify all blocks on chain and transactions and hashes and order and nonces
	for i := 1; i < len(n.blockchain); i++ {
		currBlock := n.blockchain[i]
		verified, err := currBlock.verifyBlock()
		if err != nil {
			log.Println(err)
			return false
		}
		if !verified {
			return false
		}
	}

	return true
}

func (n *node) getPublicKey() *rsa.PublicKey {
	return &n.keyPair.PublicKey
}

// addBlock to run in parallel which will run verification and then add block to blockChain of node
func (n *node) addBlock(b block, sender string) {
	// n.blockchainLock.Lock()
	n.Lock()
	defer n.Unlock()

	// After doing sufficient verification add the block to own chain
	verified, err := b.verifyBlock()
	if err != nil {
		log.Panic("Block verification failed when adding block to chain - ", err)
	}
	if !verified {
		log.Panic("unknown error occured while verifying block")
	}

	txList := b.transactionTree.getTxList()

	// Check prev hash matching to current blockChain
	if b.prevHash != n.blockchain[len(n.blockchain)-1].hash {
		// log.Println()
		log.Panicln("Node:", n.nodeID, "Sender:", sender, "Prev Hash not matching - write complete code.", "New Block PrevHash:", b.prevHash[:15], "Blockchain Last Hash:", n.blockchain[len(n.blockchain)-1].hash[:15])
	}

	// Check Balance validities
	tempBalances := createBalanceCopy(n.balances)

	for i, tx := range txList {
		tempBalances[tx.receiverID] += tx.amount
		if tx.senderID != "" {
			tempBalances[tx.senderID] -= tx.amount

			if tempBalances[tx.senderID] < 0 {
				log.Panic("Someone tried negative transactions.")
			}

		} else if i != len(txList)-1 {
			log.Panic("Someone tried to get free money")
		}
	}

	// ? Need of checking timestamp before adding to chain
	n.balances = tempBalances

	n.blockchain = append(n.blockchain, b)
	// log.Println("Node:", n.nodeID, "added new block from sender:", sender, "with hash", n.blockchain[len(n.blockchain)-1].hash, len(n.blockchain))
	log.Println("Node:", n.nodeID, "added new block from node:", sender, "PrevHash:", n.blockchain[len(n.blockchain)-1].prevHash[:15], "NewHash:", n.blockchain[len(n.blockchain)-1].hash[:15], len(n.blockchain))

	// Update selfBal if required
	// newBal := calcBalance(&b.transactionTree, n.nodeID)
	// n.selfBal += newBal
	// n.blockchainLock.Unlock()
}

func (n *node) startNode() {

	n.balances = make(map[string]float64)

	// time.Sleep(10 * time.Second)

	for !startMining {
		time.Sleep(100 * time.Millisecond)
	}

	// Start mining
	// log.Println("Starting miner for node:", n.nodeID)
	go n.startMiner()
	// All others tasks can be run in parallel

	// // * Controller can give commands to this go routine or new blocks can be discovered.

	// for true {
	// 	select {
	// 	case msg1 := <-n.cmdChannel:
	// 		log.Println(msg1)
	// 	case newBlock := <-n.receiveChannel:
	// 		log.Println(newBlock)
	// 	}

	// 	if n.nodeType == smartSingleUse || n.nodeType == smartRepeatedUse {
	// 		if n.selfBal >= n.smartBalTarget {
	// 			// TODO: Add required transaction here
	// 		}

	// 		if n.nodeType == smartSingleUse {
	// 			n.nodeType = normal
	// 		}
	// 	}
	// }

}

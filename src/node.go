package main

import (
	"crypto/rsa"
	"fmt"
	"log"
)

// Node type
const (
	normal           = iota
	malicious        = iota
	smartSingleUse   = iota
	smartRepeatedUse = iota
)

type node struct {
	nodeID                  string
	receiveChannel          chan block
	cmdChannel              chan msg
	blockchain              []block
	pendingTransactions     []transaction // to be pushed onto next block
	keyPair                 *rsa.PrivateKey
	selfBal, smartBalTarget float64
	nodeType                int
}

func (n *node) mineCoin() {
	// Create a new transaction giving yourself reward and add to pending transactions
	tx, err := createTransaction("", n.nodeID, MiningPrize)
	if err != nil {
		log.Fatal("Error mining coin - creating transaction -", err)
	}
	n.pendingTransactions = append(n.pendingTransactions, tx)

	// Create block with all pending transactions and mine to get Nonce
	prevBlockHash := n.blockchain[len(n.blockchain)-1].createHash()
	block := createBlock(n.pendingTransactions, prevBlockHash)
	block.mine()

	for _, tx := range n.pendingTransactions {
		n.selfBal += extractBalTransaction(&tx, n.nodeID)
	}

	// Push block onto chain
	n.blockchain = append(n.blockchain, block)
	n.pendingTransactions = []transaction{}

	// Broadcast new block to all nodes
	broadcastBlock(block, n.nodeID)
}

func (n *node) addTransaction(tx transaction) {
	// Add transaction to pending Transactions list after checking validity and possibility and existence etc.
	if tx.senderID == "" || tx.receiverID == "" {
		log.Fatal("Empty sender or receiver ID")
	}

	txverified, err := tx.verifyTransaction()
	if err != nil {
		log.Fatal("Error in verifiying transaction - ", err)
	}
	if !txverified {
		log.Fatal("Transaction verification failed")
	}

	if tx.amount <= 0 {
		log.Fatal("Amount <= 0")
	}

	if n.getBalance(tx.senderID) < tx.amount {
		log.Fatal("Insufficient balance for transaction")
	}

	n.pendingTransactions = append(n.pendingTransactions, tx)
}

func (n *node) getBalance(nodeID string) float64 {
	if nodeID == n.nodeID {
		return n.selfBal
	}

	bal := 0.0
	// Go over all blocks and all transactions to generate final balance of a node
	for _, b := range n.blockchain {
		bal += calcBalance(&b.transactionTree, nodeID)
	}

	return bal
}

func (n *node) verifyChain() bool {
	// ? Add check of genesis block

	// Verify all blocks on chain and transactions and hashes and order and nonces
	for i := 1; i < len(n.blockchain); i++ {
		currBlock := n.blockchain[i]
		verified, err := currBlock.verifyBlock()
		if err != nil {
			fmt.Println(err)
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

func (n *node) addBlock(b block) {
	// After doing sufficient verification add the block to own chain
	verified, err := b.verifyBlock()
	if err != nil {
		log.Fatal("Block verification failed when adding block to chain")
	}
	if !verified {
		log.Fatal("unknown error occured while verifying block")
	}

	// ? Need of checking timestamp before adding to chain
	n.blockchain = append(n.blockchain, b)

	// Update selfBal if required
	newBal := calcBalance(&b.transactionTree, n.nodeID)
	n.selfBal += newBal
}

func (n *node) startNode() {
	// * Controller can give commands to this go routine or new blocks can be discovered.

	for true {
		select {
		case msg1 := <-n.cmdChannel:
			fmt.Println(msg1)
		case newBlock := <-n.receiveChannel:
			fmt.Println(newBlock)
		}

		if n.nodeType == smartSingleUse || n.nodeType == smartRepeatedUse {
			if n.selfBal >= n.smartBalTarget {
				// TODO: Add required transaction here
			}

			if n.nodeType == smartSingleUse {
				n.nodeType = normal
			}
		}
	}

}

package main

import (
	"crypto/rsa"
	"log"
)

type node struct {
	nodeID              string
	receiveChannel      chan block
	blockchain          []block
	pendingTransactions []transaction // to be pushed onto next block
	keyPair             *rsa.PrivateKey
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

func (n *node) getBalance(nodeID string) float32 {
	//TODO: Go over all blocks and all transactions to generate final balance of a node
	return 0.0
}

func (n *node) verifyChain() bool {
	// ? Add check of genesis block

	// Verify all blocks on chain and transactions and hashes and order and nonces
	for i := 1; i < len(n.blockchain); i++ {
		currBlock := n.blockchain[i]
		if !currBlock.verifyTransactions() {
			return false
		}

		if currBlock.createHash() != currBlock.hash {
			return false
		}
	}

	return true
}

func (n *node) getPublicKey() *rsa.PublicKey {
	return &n.keyPair.PublicKey
}

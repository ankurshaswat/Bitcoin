package main

import (
	"crypto/rsa"
	"log"
	"sync"
	"time"
)

// * Node type
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
	// * Transactions to be pushed onto next block
	pendingTransactions []transaction
	keyPair             *rsa.PrivateKey
	smartBalTarget      float64
	smartNodeTarget     string
	nodeType            int
	balances            map[string]float64
}

// * Creates a transaction executing the smart contract coded into the node.
func (n *node) smartContractExecute() {

	tx, err := createTransaction(n.nodeID, n.smartNodeTarget, n.smartBalTarget)

	if err != nil {
		log.Panicln("Error creating transaction for smart contract - ", err)
	}

	tx.signTransaction(n.keyPair)

	n.balances[tx.senderID] -= tx.amount
	n.balances[tx.receiverID] += tx.amount

	broadcastTransaction(tx, n.nodeID)

	if n.nodeType == smartSingleUse {
		n.nodeType = normal
	}
}

// * Tries to broadcast a new published block.
func (n *node) tryPublishBlock() {
	// * Create a new transaction giving yourself reward and add to pending transactions
	tx, err := createTransaction("", n.nodeID, miningPrize)
	if err != nil {
		log.Panic("Error mining coin - creating transaction -", err)
	}
	n.pendingTransactions = append(n.pendingTransactions, tx)

	// * Create block with all pending transactions and mine to get Nonce
	prevBlockHash := n.blockchain[len(n.blockchain)-1].hash
	b := createBlock(n.pendingTransactions, prevBlockHash)
	mineSuccess := b.mineSingleTry()

	if !mineSuccess {
		// * Remove last transaction from pendingTransactions
		n.pendingTransactions = n.pendingTransactions[:len(n.pendingTransactions)-1]
	} else {
		for _, tx := range n.pendingTransactions {
			n.balances[tx.receiverID] += tx.amount
			if tx.senderID != "" {
				n.balances[tx.senderID] -= tx.amount
			}
		}

		// * Push block onto chain
		n.blockchain = append(n.blockchain, b)
		n.pendingTransactions = []transaction{}

		// * Broadcast new block to all nodes
		broadcastBlock(b, n.nodeID)
	}

}

// * startMiner starts an endless loop which keeps on searching
// * for a new nonce and keeps broadcasting new blocks and also checking for smart contract conditions.
// * Uses node resource lock
func (n *node) startMiner() {

	for true {
		n.Lock()
		if (n.nodeType == smartSingleUse || n.nodeType == smartRepeatedUse) && n.balances[n.nodeID] >= n.smartBalTarget {
			n.smartContractExecute()
		} else {
			n.tryPublishBlock()
		}
		n.Unlock()
	}
}

// * getPublicKey can be used as an RPC to get public key from node.
func (n *node) getPublicKey() *rsa.PublicKey {
	return &n.keyPair.PublicKey
}

// * addBlock to run in parallel which will run verification and
// * then add block to blockChain of node.
// * Uses node resource lock
func (n *node) addBlock(b block, sender string) {
	n.Lock()
	defer n.Unlock()

	// * After doing sufficient verification add the block to own chain
	verified, err := b.verifyBlock()
	if err != nil {
		log.Panic("Block verification failed when adding block to chain - ", err)
	}
	if !verified {
		log.Panic("Unknown error occured while verifying block")
	}

	txList := b.transactionTree.getTxList()

	// * Check prev hash matching to current blockChain
	if b.prevHash != n.blockchain[len(n.blockchain)-1].hash {
		log.Panicln("Node:", n.nodeID, "Sender:", sender, "Prev Hash not matching - write complete code.", "New Block PrevHash:", b.prevHash[:15], "Blockchain Last Hash:", n.blockchain[len(n.blockchain)-1].hash[:15])
	}

	// * Check Balance validities (non negative balances always)
	tempBalances := createBalanceCopy(n.balances)

	for i, tx := range txList {
		// * if tx in pendingTransactions skip check
		// * last element check is not required
		if !(i == len(txList)-1 || n.inPendingTransactions(tx)) {

			if tx.senderID == "" {
				log.Panic("Someone tried to get free money")
			} else {
				log.Panicln("Unbroadcasted transaction found. What to do here?")

				// tempBalances[tx.receiverID] += tx.amount
				// tempBalances[tx.senderID] -= tx.amount
				// if tempBalances[tx.senderID] < 0 {
				// 	log.Panic("Someone tried negative transactions.")
				// }
			}
		}
	}

	// * after checking all clear matching ones from pending transaction
	for i, tx := range txList {
		if i != len(txList)-1 {
			removed := n.removePendingTx(tx)
			if !removed {
				log.Panicln("Transaction not in pending transactions")
			}
		}
	}

	lastTx := txList[len(txList)-1]

	if lastTx.senderID != "" {
		log.Panicln("Non empty last tx in block which should have been mining prize")
	}

	if lastTx.amount != miningPrize {
		log.Panicln("Someone trying to get extra mining prize")
	}

	tempBalances[lastTx.receiverID] += lastTx.amount

	// ? Need of checking timestamp before adding to chain

	n.balances = tempBalances
	n.blockchain = append(n.blockchain, b)

	log.Println("Node:", n.nodeID, "added new block from node:", sender, "PrevHash:", n.blockchain[len(n.blockchain)-1].prevHash[:15], "NewHash:", n.blockchain[len(n.blockchain)-1].hash[:15], len(n.blockchain))
}

// * startNode spins off the mining go routine.
func (n *node) startNode() {
	n.balances = make(map[string]float64)

	for !startMining {
		time.Sleep(100 * time.Millisecond)
	}

	// * Start mining
	go n.startMiner()

	// * All others tasks can be run in parallel as RPCs.
}

// * receiveTransaction receives a broadcast of a transaction
// * then adds it to pendingTransactions after verification.
func (n *node) receiveTransaction(tx transaction) {

	//* Add transaction to pending Transactions list after checking validity and possibility and existence etc.

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

	n.balances[tx.senderID] -= tx.amount
	n.balances[tx.receiverID] += tx.amount

	n.pendingTransactions = append(n.pendingTransactions, tx)
}

func (n *node) inPendingTransactions(tx transaction) bool {
	for _, txLocal := range n.pendingTransactions {
		if txLocal.hash == tx.hash {
			return true
		}
	}
	return false
}

func (n *node) removePendingTx(tx transaction) bool {
	pos := -1
	for i, txLocal := range n.pendingTransactions {
		if txLocal.hash == tx.hash {
			pos = i
			break
		}
	}

	if pos == -1 {
		return false

	}

	n.pendingTransactions = append(n.pendingTransactions[:pos], n.pendingTransactions[pos+1:]...)
	return true

}

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

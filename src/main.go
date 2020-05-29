package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

// RSAbitSize ... Bit size for key pair
const (
	RSAbitSize      = 10
	MiningDificulty = 2
	MiningPrize     = 1
)

func generateSHA256Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}

func generateKeyPair() *rsa.PrivateKey {
	reader := rand.Reader
	key, _ := rsa.GenerateKey(reader, RSAbitSize)
	return key
}

type transaction struct {
	senderID, receiverID, signature string
	amount                          float32
	timestamp                       time.Time
}

func (t *transaction) getHash() string {
	amountString := fmt.Sprintf("%f", t.amount)
	s := t.senderID + t.receiverID + amountString
	hash := generateSHA256Hash(s)
	return hash
}

func (t *transaction) signTransaction(keyPair *rsa.PrivateKey) {
	// check if public key is senders public key (probab write a function for node to get its public key)
	// gethash of transaction and sign it with key
	// return after saving signature in transaction
}

func (t *transaction) verifyTransaction() bool {
	// if (no from address) {
	//  No from address means mining reward self added (how to check if valid - ans: using proof-of-work)
	// return true
	// }

	// if (no signature) {
	// no signature - invalid transaction
	// return false
	// }

	// get public key of from address
	// verify signature with hash of transaction

	return false
}

func createTransaction(senderID string, receiverID string, amount float32) transaction {
	return transaction{senderID: senderID, receiverID: receiverID, amount: amount, timestamp: time.Now()}
}

type treeBlockLeaf struct {
	hash        string
	left, right transaction
}
type treeBlock struct {
	hash        string
	left, right treeBlockLeaf
}

type merkleTree struct {
	root treeBlock
}

func createMerkleTree(tList []transaction) {
	// super duper easy to write
}

type block struct {
	timestamp, nonce int
	prevHash, string string
	transactions     merkleTree
}

func (b *block) createHash() string {
	// use prevHash,nonce,timestamp,merkleTree root hash
	// super duper easy again
	return ""
}

func (b *block) mine() int {
	// find a nonce (by incrementing and checking) so that sha256 has difficulty number of zeros
	return 0
}

func (b *block) verifyBlock() bool {

	// Check all transactions
	// Check merkle tree???
	// Check nonce

	return false
}

type node struct {
	nodeID              string
	receiveChannel      chan block
	blockchain          []block
	pendingTransactions []transaction // to be pushed onto next block
}

func (n *node) mineCoin() {
	// Create a new transaction giving yourself reward and add to pending transactions
	// Create block with all pending transactions and mine to get Nonce
	// Push block onto chain
	// Broadcast new block
}

func (n *node) addTransaction() {
	// Add transaction to pending Transactions list after checking validity and possibility and existence etc.
}

func (n *node) getBalance(nodeID string) float32 {
	// Go over all blocks and all transactions to generate final balance of a node
	return 0.0
}

func (n *node) verifyChain() bool {
	// Verify all blocks on chain and transactions and hashes and order and nonces

	return false
}

func main() {
	// Create genesis block
	// Create new nodes in loop (while also providing genesis block ref)
	// Order a node to do a transaction
}

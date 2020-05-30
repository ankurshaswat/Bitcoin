package main

import (
	"strconv"
	"strings"
	"time"
)

type block struct {
	timestamp       time.Time
	nonce           int
	prevHash, hash  string
	transactionTree treeBlock
}

func (b *block) createHashCustom(nonce int) string {
	// use prevHash,nonce,timestamp,merkleTree root hash
	return generateSHA256Hash(b.prevHash + strconv.Itoa(nonce) + strconv.FormatInt(b.timestamp.UnixNano(), 10) + b.transactionTree.hash)
}

func (b *block) createHash() string {
	return b.createHashCustom(b.nonce)
}

func (b *block) checkProofOfWork(hash string) bool {
	subString := hash[0:MiningDificulty]
	correct := strings.Repeat("0", MiningDificulty)
	return subString == correct
}

func (b *block) mine() {
	// find a nonce (by incrementing and checking) so that sha256 has difficulty number of zeros
	nonce := b.nonce + 1
	hash := b.createHashCustom(nonce)

	for !b.checkProofOfWork(hash) {
		nonce++
		hash = b.createHashCustom(nonce)
	}

	b.hash = hash
	b.nonce = nonce
}

func (b *block) verifyTransactions() bool {

	// Check merkle tree with all transactions checked recursively
	if !b.transactionTree.verifyTree() {
		return false
	}

	// ? Check anything else at block level??? like nonce

	return true
}

func createBlock(transactions []transaction, prevHash string) block {
	transactionTree := createMerkleTree(transactions)
	return block{timestamp: time.Now(), prevHash: prevHash, transactionTree: transactionTree}
}

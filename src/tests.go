package main

import (
	"log"
)

func runTests() {
	originalDebug := debugP
	debugP = true

	txTests()
	treeTests()
	blockTests()
	nodeTests()

	debugP = originalDebug
}

func txTests() {
	log.Print("\n-----Running tests for transaction module-----\n\n")

	tx, err := createTransaction("0", "1", -1)
	if err != nil {
		log.Println("✅ Successfully caught negative amount error -", err)
	} else {
		log.Println("❌ Failed to catch error")
	}

	tx, err = createTransaction("0", "1", 0)
	if err != nil {
		log.Println("✅ Successfully caught 0 amount error -", err)
	} else {
		log.Println("❌ Failed to catch error")
	}

	tx, err = createTransaction("0", "0", 10)
	if err != nil {
		log.Println("✅ Successfully caught same id error -", err)
	} else {
		log.Println("❌ Failed to catch error")
	}

	tx, err = createTransaction("0", "1", 100)
	if err == nil {
		log.Println("✅ Successfully created transaction.")
	} else {
		log.Println("❌ Unknow error encountenered creating transaction - ", err)
	}

	log.Println("Testing getHash")
	hash := tx.getHash()
	log.Printf("✅ Got hash %v \n", hash)

	log.Println("Testing verifyTransaction")
	verified, err := tx.verifyTransaction()
	if err != nil {
		log.Println("✅ Successfully caught unsigned transaction -", err)
	} else {
		log.Println("❌ Failed to catch error")
	}

	log.Println("Trying to generate key pair")
	keyPair := generateKeyPair()
	log.Println("Key generation working -", keyPair)

	log.Println("Creating node for signing tests")
	// n :=
	nodeList = append(nodeList, createNode("0", []block{}, normal))

	log.Println("Trying to sign transaction")
	tx.signTransaction(nodeList[0].keyPair)

	log.Println("Testing verifyTransaction")
	verified, err = tx.verifyTransaction()
	if err == nil {
		log.Println("✅ Successfully signed transaction")
	} else {
		log.Println("❌ Error in signing transaction - ", err)
	}

	log.Println(verified)
}

func treeTests() {
	log.Print("\n-----Running tests for merkle tree module-----\n\n")

	t1, err := createTransaction("0", "1", 100)
	if err != nil {
		log.Println(err)
	}
	t2, err := createTransaction("1", "2", 50)
	t3, err := createTransaction("2", "0", 10)

	tList := []transaction{t1, t2, t3}
	tree := createMerkleTree(tList)
	log.Println("Merkle tree created")
	log.Println(tree)

	log.Println("Checking verify tree function")
	res, err := tree.verifyTree()
	if err != nil {
		log.Println("✅ Successfully caught unsigned transaction in tree - ", err)
	} else {
		log.Println("❌ Failed to catch unverified transaction in tree")
	}
	log.Println(res)
}

func blockTests() {
	log.Print("\n-----Running tests for block module-----\n\n")

	t1, err := createTransaction("0", "1", 100)
	if err != nil {
		log.Println(err)
	}
	t2, err := createTransaction("1", "2", 50)
	t3, err := createTransaction("2", "0", 10)

	tList := []transaction{t1, t2, t3}

	b := createBlock(tList, "")
	log.Println("Block succesfully created")
	log.Println(b)

	log.Println("Test createHashCustom")
	s := b.createHashCustom(0)
	log.Println("CreateHashCustom result - ", s)

	log.Println("Test createHash")
	s = b.createHash()
	log.Println("CreateHash result - ", s)

	log.Println("Check proof of work function")
	res := b.checkProofOfWork(s)
	log.Println("Result of proof of work check with last hash - ", res)

	log.Println("Check mine function")
	b.mine()
	log.Printf("Result after mining - hash:%v nonce:%v\n", b.hash, b.nonce)

	log.Println("Checking verify transactions")
	verified, err := b.verifyTransactions()
	if err != nil {
		log.Println("✅ Successfully caught unsigned transaction in block - ", err)
	} else {
		log.Println("❌ Failed to catch unverified transaction in block")
	}
	log.Println(verified)
}

func nodeTests() {
	log.Print("\n-----Running tests for node module-----\n\n")

	// Handle genesis block creation before writing rest of the tests
	log.Println("Creating genesis block")
	testGenBlock := createBlock([]transaction{}, "")

	nodeID := "0"
	// receiveChannel := make(chan block, 10)
	blockChain := []block{testGenBlock}

	log.Println("Creating Normal Node")
	n := createNode(nodeID, blockChain, normal)
	log.Println("node created - ", &n)

	log.Println("Mining coin")
	n.startMiner()
	log.Println("Coin mined")

	log.Println("Adding Transaction")
	tx, err := createTransaction("0", "1", 100)
	if err != nil {
		log.Println(err)
	}

	log.Println("Trying to add unsigned transaction")
	n.addTransaction(tx)

	// log.Println("Getting public key of node")
	// pbKey := n.getPublicKey()
	// log.Println("Public key of node - ", pbKey)

	// log.Println("Verifying chain of blocks of node")

}

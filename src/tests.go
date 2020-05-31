package main

import "fmt"

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
	fmt.Print("\n-----Running tests for transaction module-----\n\n")

	tx, err := createTransaction("0", "1", -1)
	if err != nil {
		fmt.Println("✅ Successfully caught negative amount error -", err)
	} else {
		fmt.Println("❌ Failed to catch error")
	}

	tx, err = createTransaction("0", "1", 0)
	if err != nil {
		fmt.Println("✅ Successfully caught 0 amount error -", err)
	} else {
		fmt.Println("❌ Failed to catch error")
	}

	tx, err = createTransaction("0", "0", 10)
	if err != nil {
		fmt.Println("✅ Successfully caught same id error -", err)
	} else {
		fmt.Println("❌ Failed to catch error")
	}

	tx, err = createTransaction("0", "1", 100)
	if err == nil {
		fmt.Println("✅ Successfully created transaction.")
	} else {
		fmt.Println("❌ Unknow error encountenered creating transaction - ", err)
	}

	fmt.Println("Testing getHash")
	hash := tx.getHash()
	fmt.Printf("✅ Got hash %v \n", hash)

	fmt.Println("Testing verifyTransaction")
	verified, err := tx.verifyTransaction()
	if err != nil {
		fmt.Println("✅ Successfully caught unsigned transaction -", err)
	} else {
		fmt.Println("❌ Failed to catch error")
	}

	// fmt.Println("Trying to generate key pair")
	// keyPair := generateKeyPair()

	// fmt.Println("Trying to sign transaction")
	// tx.signTransaction(keyPair)

	// fmt.Println("Testing verifyTransaction")
	// verified, err = tx.verifyTransaction()
	// if err == nil {
	// 	fmt.Println("✅ Successfully signed transaction")
	// } else {
	// 	fmt.Println("❌ Error in signing transaction - ", err)
	// }

	fmt.Println(verified)
}

func treeTests() {
	fmt.Print("\n-----Running tests for merkle tree module-----\n\n")

	t1, err := createTransaction("0", "1", 100)
	if err != nil {
		fmt.Println(err)
	}
	t2, err := createTransaction("1", "2", 50)
	t3, err := createTransaction("2", "0", 10)

	tList := []transaction{t1, t2, t3}
	tree := createMerkleTree(tList)
	fmt.Println("Merkle tree created")
	fmt.Println(tree)

	fmt.Println("Checking verify tree function")
	res, err := tree.verifyTree()
	if err != nil {
		fmt.Println("✅ Successfully caught unsigned transaction in tree - ", err)
	} else {
		fmt.Println("❌ Failed to catch unverified transaction in tree")
	}
	fmt.Println(res)
}

func blockTests() {
	fmt.Print("\n-----Running tests for block module-----\n\n")

	t1, err := createTransaction("0", "1", 100)
	if err != nil {
		fmt.Println(err)
	}
	t2, err := createTransaction("1", "2", 50)
	t3, err := createTransaction("2", "0", 10)

	tList := []transaction{t1, t2, t3}

	b := createBlock(tList, "")
	fmt.Println("Block succesfully created")
	fmt.Println(b)

	fmt.Println("Test createHashCustom")
	s := b.createHashCustom(0)
	fmt.Println("CreateHashCustom result - ", s)

	fmt.Println("Test createHash")
	s = b.createHash()
	fmt.Println("CreateHash result - ", s)

	fmt.Println("Check proof of work function")
	res := b.checkProofOfWork(s)
	fmt.Println("Result of proof of work check with last hash - ", res)

	fmt.Println("Check mine function")
	b.mine()
	fmt.Printf("Result after mining - hash:%v nonce:%v\n", b.hash, b.nonce)

	fmt.Println("Checking verify transactions")
	verified, err := b.verifyTransactions()
	if err != nil {
		fmt.Println("✅ Successfully caught unsigned transaction in block - ", err)
	} else {
		fmt.Println("❌ Failed to catch unverified transaction in block")
	}
	fmt.Println(verified)
}

func nodeTests() {
	fmt.Print("\n-----Running tests for node module-----\n\n")

	// TODO: Handle genesis block creation before writing rest of the tests

	nodeID := "0"
	receiveChannel := make(chan block, 10)
	blockChain := []block{}

	fmt.Println("Creating Node")
	n := createNode(nodeID, receiveChannel, blockChain)
	fmt.Println("node created - ", n)

	fmt.Println("Mining coin")
	n.mineCoin()
	fmt.Println("Coin mined")

	fmt.Println("Adding Transaction")
	tx, err := createTransaction("0", "1", 100)
	if err != nil {
		fmt.Println(err)
	}
	n.addTransaction(tx)

	fmt.Println("Getting public key of node")
	pbKey := n.getPublicKey()
	fmt.Println("Public key of node - ", pbKey)

	// fmt.Println("Verifying chain of blocks of node")

}

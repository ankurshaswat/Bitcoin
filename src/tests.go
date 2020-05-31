package main

import "fmt"

func runTests() {
	originalDebug := debugP
	debugP = true

	txTests()

	debugP = originalDebug
}

func txTests() {
	fmt.Println("Running tests for transaction module")

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

	fmt.Println("Trying to generate key pair")
	keyPair := generateKeyPair()

	fmt.Println("Trying to sign transaction")
	tx.signTransaction(keyPair)

	// fmt.Println("Testing verifyTransaction")
	// verified, err = tx.verifyTransaction()
	// if err == nil {
	// 	fmt.Println("✅ Successfully signed transaction")
	// } else {
	// 	fmt.Println("❌ Error in signing transaction - ", err)
	// }

	// fmt.Println(verified)
}

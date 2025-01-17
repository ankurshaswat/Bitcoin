package main

import "fmt"

// type treeBlockLeaf struct {
// 	hash        string
// 	left, right *transaction
// }

type merkleTree struct {
	hash          string
	leaf          bool
	left, right   *merkleTree
	leftT, rightT *transaction
}

// type merkleTree struct {
// 	hash string
// 	left,right *merkleTree
// }

func (tree *merkleTree) verifyTree() (bool, error) {

	if tree.leaf {
		verified, err := tree.leftT.verifyTransaction()
		stringToHash := tree.leftT.hash

		if err != nil {
			return false, err
		}

		if tree.rightT != nil {
			verify2, err := tree.rightT.verifyTransaction()
			stringToHash = stringToHash + tree.rightT.hash
			if err != nil {
				return false, err
			}
			verified = verified && verify2
		}

		if !verified {
			return false, fmt.Errorf("Sub tree transaction verification failed")
		}

		if tree.hash != generateSHA256Hash(stringToHash) {
			return false, fmt.Errorf("leaf hash matching failed")
		}

	} else {
		verified, err := tree.left.verifyTree()
		stringToHash := tree.left.hash

		if err != nil {
			return false, err
		}

		if tree.right != nil {
			verify2, err := tree.right.verifyTree()
			stringToHash = stringToHash + tree.right.hash
			if err != nil {
				return false, err
			}
			verified = verified && verify2
		}

		if !verified {
			return false, fmt.Errorf("Sub tree verification failed")
		}

		if tree.hash != generateSHA256Hash(stringToHash) {
			return false, fmt.Errorf("sub tree hash matching failed")
		}
	}

	return true, nil
}

func (tree *merkleTree) getTxList() []transaction {
	txList := []transaction{}

	if !tree.leaf {
		txList = tree.left.getTxList()
		if tree.right != nil {
			txList = append(txList, tree.right.getTxList()...)
		}
	} else {
		txList = append(txList, *tree.leftT)
		if tree.rightT != nil {
			txList = append(txList, *tree.rightT)
		}
	}

	return txList
}

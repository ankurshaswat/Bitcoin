package main

import "fmt"

// type treeBlockLeaf struct {
// 	hash        string
// 	left, right *transaction
// }

type treeBlock struct {
	hash          string
	leaf          bool
	left, right   *treeBlock
	leftT, rightT *transaction
}

// type merkleTree struct {
// 	hash string
// 	left,right *merkleTree
// }

func (tb *treeBlock) verifyTree() (bool, error) {

	if tb.leaf {
		verified, err := tb.leftT.verifyTransaction()
		stringToHash := tb.leftT.getHash()

		if err != nil {
			return false, err
		}

		if tb.rightT != nil {
			verify2, err := tb.rightT.verifyTransaction()
			stringToHash = stringToHash + tb.rightT.getHash()
			if err != nil {
				return false, err
			}
			verified = verified && verify2
		}

		if !verified {
			return false, fmt.Errorf("Sub tree transaction verification failed")
		}

		if tb.hash != generateSHA256Hash(stringToHash) {
			return false, fmt.Errorf("leaf hash matching failed")
		}

	} else {
		verified, err := tb.left.verifyTree()
		stringToHash := tb.left.hash

		if err != nil {
			return false, err
		}

		if tb.right != nil {
			verify2, err := tb.right.verifyTree()
			stringToHash = stringToHash + tb.right.hash
			if err != nil {
				return false, err
			}
			verified = verified && verify2
		}

		if !verified {
			return false, fmt.Errorf("Sub tree verification failed")
		}

		if tb.hash != generateSHA256Hash(stringToHash) {
			return false, fmt.Errorf("sub tree hash matching failed")
		}
	}

	return true, nil
}

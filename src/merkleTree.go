package main

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

func createMerkleTree(tList []transaction) treeBlock {

	blockList := []treeBlock{}

	for i := 0; i < len(tList); i += 2 {
		trans1 := tList[i]
		if i+1 == len(tList) {
			// If only one left
			hash := generateSHA256Hash(trans1.getHash())
			newLeafBlock := treeBlock{leaf: true, leftT: &trans1, hash: hash}
			blockList = append(blockList, newLeafBlock)
		} else {
			// if more than one available
			trans2 := tList[i+1]
			hash := generateSHA256Hash(trans1.getHash() + trans2.getHash())
			newLeafBlock := treeBlock{leaf: true, leftT: &trans1, rightT: &trans2, hash: hash}
			blockList = append(blockList, newLeafBlock)
		}
	}

	for len(blockList) > 1 {
		newBlockList := []treeBlock{}

		for i := 0; i < len(blockList); i += 2 {
			block1 := blockList[i]
			if i+1 == len(blockList) {
				// If only one left
				hash := generateSHA256Hash(block1.hash)
				newTreeBlock := treeBlock{leaf: false, left: &block1, hash: hash}
				newBlockList = append(newBlockList, newTreeBlock)
			} else {
				// if two available
				block2 := blockList[i+1]
				hash := generateSHA256Hash(block1.hash + block2.hash)
				newTreeBlock := treeBlock{leaf: false, left: &block1, right: &block2, hash: hash}
				newBlockList = append(newBlockList, newTreeBlock)
			}
		}
		blockList = newBlockList
	}

	return blockList[0]
}

func (tb *treeBlock) verifyTree() bool {

	if tb.leaf {

		if tb.rightT != nil {

			if !tb.leftT.verifyTransaction() || !tb.rightT.verifyTransaction() {
				return false
			}

			if tb.hash != generateSHA256Hash(tb.leftT.getHash()+tb.rightT.getHash()) {
				return false
			}
		} else {
			if !tb.leftT.verifyTransaction() {
				return false
			}
			if tb.hash != generateSHA256Hash(tb.leftT.getHash()) {
				return false
			}
		}

	} else {
		if tb.right != nil {
			if !tb.left.verifyTree() || !tb.right.verifyTree() {
				return false
			}

			if tb.hash != generateSHA256Hash(tb.left.hash+tb.right.hash) {
				return false
			}
		} else {
			if !tb.left.verifyTree() {
				return false
			}

			if tb.hash != generateSHA256Hash(tb.left.hash) {
				return false
			}
		}
	}

	return true
}

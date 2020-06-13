package main

import (
	"fmt"
	"time"
)

func createTransaction(senderID string, receiverID string, amount float64) (transaction, error) {
	debug(fmt.Sprintf("Creating transaction - senderID:%v receiverID:%v amount:%v", senderID, receiverID, amount))

	if amount <= 0 {
		return transaction{}, fmt.Errorf("Wrong Amount Passed %v", amount)
	}

	if receiverID == "" {
		return transaction{}, fmt.Errorf("Empty receiver ID")
	}

	if senderID == receiverID {
		return transaction{}, fmt.Errorf("Sender and receiver id cannot be same %v", senderID)
	}

	return transaction{senderID: senderID, receiverID: receiverID, amount: amount, timestamp: time.Now()}, nil
}

func broadcastTransaction(t transaction, selfID string) {
	for i := 0; i < len(nodeList); i++ {
		nodeInstance := &nodeList[i]
		if nodeInstance.nodeID != selfID {
			// txMsg := msg{msgType: transactionBroadcast, tx: t}
			// nodeInstance.cmdChannel <- txMsg
		}
	}
}

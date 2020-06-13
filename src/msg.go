package main

// Msg Type enum
const (
	transactionBroadcast = iota
	// mine     = iota
)

type msg struct {
	msgType int
	tx      transaction
}

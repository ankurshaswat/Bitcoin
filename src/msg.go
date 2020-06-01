package main

// Msg Type enum
const (
	mine     = iota
	transact = iota
)

type msg struct {
	msgType int
}

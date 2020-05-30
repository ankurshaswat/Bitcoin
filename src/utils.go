package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
)

// RSAbitSize ... Bit size for key pair
const (
	RSAbitSize      = 10
	MiningDificulty = 2
	MiningPrize     = 1
)

func generateSHA256Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}

func generateKeyPair() *rsa.PrivateKey {
	reader := rand.Reader
	key, _ := rsa.GenerateKey(reader, RSAbitSize)
	return key
}

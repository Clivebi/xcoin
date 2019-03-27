package proxy

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"
)

var publickkeys []crypto.PublicKey
var privatekeys []crypto.PrivateKey
var userids []string

const (
	rootUser    = 0
	normalUserA = 1
	normalUserB = 2
	apiURL      = "http://127.0.0.1:8789/callapi.do"
)

func initEnv() {
	publickkeys = make([]crypto.PublicKey, 3)
	privatekeys = make([]crypto.PrivateKey, 3)
	userids = make([]string, 3)

	for i := 0; i < 3; i++ {
		key, _ := rsa.GenerateKey(rand.Reader, 2048)
		publickkeys[i] = key.Public()
		privatekeys[i] = key
	}
}

func addressFromPublicKey(pubkey interface{}) string {
	text, _ := PublicKeyToString(pubkey)
	hashed := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hashed[:])
}

func TestUser(t *testing.T) {
	initEnv()
	t.Log("add root")
	wallet, err := RegisgerWallet(apiURL, addressFromPublicKey(publickkeys[rootUser]),
		publickkeys[rootUser], "USD")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(wallet)

	wallet, err = RegisgerWallet(apiURL, addressFromPublicKey(publickkeys[normalUserA]),
		publickkeys[normalUserA], "USD")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(wallet)

	wallet, err = RegisgerWallet(apiURL, addressFromPublicKey(publickkeys[normalUserB]),
		publickkeys[normalUserB], "USD")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(wallet)
}

func TestSend(t *testing.T) {
	t.Log("Admin send 1000 to userA")
	wallet, err := Send(apiURL, addressFromPublicKey(publickkeys[rootUser]),
		addressFromPublicKey(publickkeys[normalUserA]), 1000, privatekeys[rootUser])
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(5 * time.Second)
	walletA, err := GetWallet(apiURL, addressFromPublicKey(publickkeys[normalUserA]))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("After send Admin Token:", wallet.Token, " userA Token:", walletA.Token)

	t.Log("userA send 500 to userB")
	wallet, err = Send(apiURL, addressFromPublicKey(publickkeys[normalUserA]),
		addressFromPublicKey(publickkeys[normalUserB]), 500, privatekeys[normalUserA])
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(5 * time.Second)
	walletA, err = GetWallet(apiURL, addressFromPublicKey(publickkeys[normalUserB]))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("After send userA Token:", wallet.Token, " userB Token:", walletA.Token)
}

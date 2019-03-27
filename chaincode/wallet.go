package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type group struct {
	Name   string   `json:"name"`
	Root   string   `json:"root"`
	Wallet []string `json:"wallet"`
}

func (o group) toBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

func (o group) save(stub shim.ChaincodeStubInterface) error {
	key := "group_item_" + o.Name
	return stub.PutState(key, o.toBuffer())
}

func getGroup(stub shim.ChaincodeStubInterface, name string) *group {
	key := "group_item_" + name
	g := &group{}
	buf, _ := stub.GetState(key)
	if buf != nil {
		if err := json.Unmarshal(buf, g); err != nil {
			return nil
		}
		return g
	}
	return nil
}

type wallet struct {
	PublicKey string  `json:"pub_key"`
	Address   string  `json:"address"`
	Token     float64 `json:"token"`
	Type      int     `json:"type"` // 0 root 1 other
	Group     string  `json:"group"`
}

func (o wallet) save(stub shim.ChaincodeStubInterface) error {
	key := "" + o.Address
	return stub.PutState(key, o.toBuffer())
}

func (o wallet) toBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

func (o wallet) checkSignature(args string, sig string) error {
	return checkSignature(args, sig, o.PublicKey)
}

func getWallet(stub shim.ChaincodeStubInterface, address string) *wallet {
	w := &wallet{}
	buf, _ := stub.GetState(address)
	if buf != nil {
		if err := json.Unmarshal(buf, w); err != nil {
			return nil
		}
		return w
	}
	return nil
}

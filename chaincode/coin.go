package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	DEBUG = true
)

type CoinChaincode struct {
}

func (t *CoinChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("init")
	NewManger(stub)
	return shim.Success([]byte("init ok"))
}

func (t *CoinChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoke")
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Invoke ", function, args[0])
	switch function {
	case "adduser":
		return t.addUser(stub, args)
	case "getuser":
		return t.getUser(stub, args)
	case "upgradeuser":
		return t.upgradeUser(stub, args)
	case "send":
		return t.sendTranscation(stub, args)
	default:
		return shim.Error("invalid function:" + function)
	}
}

func (t *CoinChaincode) checkSignature(args string, sig string, key string, stub shim.ChaincodeStubInterface) error {
	if DEBUG {
		return nil
	}
	pkbuf, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}
	pk, err := x509.ParsePKCS1PublicKey(pkbuf)
	if err != nil {
		return err
	}
	signature, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(args))
	return rsa.VerifyPKCS1v15(pk, crypto.SHA256, hashed[:], signature)
}

func (t *CoinChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	req := &AddUserOffer{}

	err := json.Unmarshal([]byte(args[0]), req)
	if err != nil {
		return shim.Error("Incorrect arguments. detail:" + err.Error())
	}

	err = t.checkSignature(args[0], args[1], req.PublickKey, stub)
	if err != nil {
		return shim.Error("check signature failed,detail:" + err.Error())
	}
	manger := NewManger(stub)
	buf, err := manger.AddUser(req.PublickKey, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(buf)
}

func (t *CoinChaincode) getUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	req := &GetUserOffser{}

	err := json.Unmarshal([]byte(args[0]), req)
	if err != nil {
		return shim.Error("Incorrect arguments. detail:" + err.Error())
	}

	manger := NewManger(stub)
	Caller := manger.getUser(req.CallID, stub)
	if Caller == nil {
		return shim.Error("Incorrect arguments. detail:" + " Get Caller failed id=" + req.CallID)
	}

	err = t.checkSignature(args[0], args[1], Caller.PubKey, stub)
	if err != nil {
		return shim.Error("check signature failed,detail:" + err.Error())
	}
	buf, err := manger.GetUser(req.UserID, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(buf)
}

func (t *CoinChaincode) upgradeUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	req := &UpgradeUserOffser{}

	err := json.Unmarshal([]byte(args[0]), req)
	if err != nil {
		return shim.Error("Incorrect arguments. detail:" + err.Error())
	}

	manger := NewManger(stub)

	Caller := manger.getUser(req.CallID, stub)
	if Caller == nil {
		return shim.Error("Incorrect arguments. detail:" + " Get Caller failed id=" + req.CallID)
	}

	err = t.checkSignature(args[0], args[1], Caller.PubKey, stub)
	if err != nil {
		return shim.Error("check signature failed,detail:" + err.Error())
	}
	err = manger.UpgradeUser(Caller, req.UserID, req.Limit, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *CoinChaincode) sendTranscation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	req := &SendTranscationOffser{}

	err := json.Unmarshal([]byte(args[0]), req)
	if err != nil {
		return shim.Error("Incorrect arguments. detail:" + err.Error())
	}
	manger := NewManger(stub)
	Caller := manger.getUser(req.CallID, stub)
	if Caller == nil {
		return shim.Error("Incorrect arguments. detail:" + " Get Caller failed id=" + req.CallID)
	}

	err = t.checkSignature(args[0], args[1], Caller.PubKey, stub)
	if err != nil {
		return shim.Error("check signature failed,detail:" + err.Error())
	}
	err = manger.Send(Caller, req.ToUser, req.Coin, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(CoinChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

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

type CoinChaincode struct {
	debug  bool
	manger *UserManger
}

func (t *CoinChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	buf, err := stub.GetState(MANGER_KEY)
	if err != nil {
		obj := &UserManger{}
		err := json.Unmarshal(buf, obj)
		if err == nil {
			t.manger = obj
		}
	}
	if t.manger == nil {
		t.manger = &UserManger{
			Users:    make([]string, 1),
			Sellers:  make([]string, 1),
			LogIndex: 0,
		}
		t.manger.save(stub)
	}
	t.debug = false
	_, args := stub.GetFunctionAndParameters()
	for _, v := range args {
		if v == "debug" {
			t.debug = true
		}
	}
	return shim.Success(nil)
}

func (t *CoinChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
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

func (t *CoinChaincode) checkSignature(args string, sig string, key string) error {
	if t.debug {
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

	err = t.checkSignature(args[0], args[1], req.PublickKey)
	if err != nil {
		return shim.Error("check signature failed,detail:" + err.Error())
	}

	buf, err := t.manger.AddUser(req.PublickKey, stub)
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

	Caller := t.manger.getUser(req.CallID, stub)
	if Caller == nil {
		return shim.Error("Incorrect arguments. detail:" + " Get Caller failed id=" + req.CallID)
	}

	err = t.checkSignature(args[0], args[1], Caller.PubKey)
	if err != nil {
		return shim.Error("check signature failed,detail:" + err.Error())
	}
	buf, err := t.manger.GetUser(req.UserID, stub)
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

	Caller := t.manger.getUser(req.CallID, stub)
	if Caller == nil {
		return shim.Error("Incorrect arguments. detail:" + " Get Caller failed id=" + req.CallID)
	}

	err = t.checkSignature(args[0], args[1], Caller.PubKey)
	if err != nil {
		return shim.Error("check signature failed,detail:" + err.Error())
	}
	err = t.manger.UpgradeUser(Caller, req.UserID, req.Limit, stub)
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

	Caller := t.manger.getUser(req.CallID, stub)
	if Caller == nil {
		return shim.Error("Incorrect arguments. detail:" + " Get Caller failed id=" + req.CallID)
	}

	err = t.checkSignature(args[0], args[1], Caller.PubKey)
	if err != nil {
		return shim.Error("check signature failed,detail:" + err.Error())
	}
	err = t.manger.Send(Caller, req.ToUser, req.Coin, stub)
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

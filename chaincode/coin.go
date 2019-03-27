package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

// CoinChaincode 实际操作类
type CoinChaincode struct {
}

//Init initlaize chaincode
func (t *CoinChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("init")
	return shim.Success([]byte("init ok"))
}

//Invoke function call interface
func (t *CoinChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var buf []byte
	var err error
	function, args := stub.GetFunctionAndParameters()
	if function != "callapi" {
		return shim.Error("chaincode function not equal :callapi")
	}
	if len(args) != 2 {
		return shim.Error("chaincode args count not equal 2")
	}

	rawreq := args[0]
	req := &Request{}
	err = json.Unmarshal([]byte(rawreq), req)
	if err != nil {
		return shim.Error(err.Error())
	}
	sig := args[1]

	switch req.Function {
	case "register", "addwallet":
		buf, err = t.registerWallet(req.Args, stub)
	case "getwallet":
		buf, err = t.getWallet(req.Args, stub)
	case "transfer", "send":
		buf, err = t.transfer(req.Args, rawreq, sig, stub)
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(buf)
}

//[walletAddress,ecdsa public key,groupname]
func (t *CoinChaincode) registerWallet(args []string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("mismatch arguments [walletAddress,ecdsa public key,groupname]")
	}
	_, err := stringToPublicKey(args[1])
	if err != nil {
		return nil, err
	}
	w := getWallet(stub, args[0])
	if w != nil {
		return nil, errors.New("address " + args[0] + " is used")
	}
	w = &wallet{
		Token:     0,
		Type:      1,
		Address:   args[0],
		PublicKey: args[1],
		Group:     args[2],
	}
	g := getGroup(stub, args[2])
	if g == nil {
		g = &group{
			Name:   args[2],
			Root:   w.Address,
			Wallet: []string{},
		}
		err = g.save(stub)
		if err != nil {
			return nil, err
		}
		w.Type = 0
	} else {
		g.Wallet = append(g.Wallet, w.Address)
		err = g.save(stub)
		if err != nil {
			return nil, err
		}
	}
	err = w.save(stub)
	if err != nil {
		return nil, err
	}
	return w.toBuffer(), nil
}

func (t *CoinChaincode) getWallet(args []string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("mismatch arguments [walletAddress]")
	}
	w := getWallet(stub, args[0])
	if w == nil {
		return nil, errors.New("address " + args[0] + " not found")
	}
	return w.toBuffer(), nil
}

//[fromwallet,towallet,ammount]
func (t *CoinChaincode) transfer(args []string, rawreq string, sig string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("mismatch arguments [fromwallet,towallet,ammount]")
	}
	from := getWallet(stub, args[0])
	if from == nil {
		return nil, errors.New("address " + args[0] + " not found")
	}
	if err := from.checkSignature(rawreq, sig); err != nil {
		return nil, err
	}
	to := getWallet(stub, args[1])
	if to == nil {
		return nil, errors.New("address " + args[1] + " not found")
	}
	if from.Group != to.Group {
		return nil, errors.New("fromwallet and towallet not in same group")
	}
	amount, err := strconv.ParseFloat(args[2], 10)
	if err != nil {
		return nil, err
	}
	if from.Type == 1 && amount > from.Token {
		return nil, errors.New("from wallet not have enough token")
	}
	if from.Type == 0 {
		from.Token += amount
	} else {
		from.Token -= amount
	}
	if to.Type == 0 {
		to.Token -= amount
	} else {
		to.Token += amount
	}
	err = from.save(stub)
	if err != nil {
		return nil, err
	}
	err = to.save(stub)
	if err != nil {
		if from.Type == 0 {
			from.Token -= amount
		} else {
			from.Token += amount
		}
		from.save(stub)
		return nil, err
	}
	return from.toBuffer(), nil
}

func main() {
	err := shim.Start(new(CoinChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

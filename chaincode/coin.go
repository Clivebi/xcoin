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
	"strconv"
)

const (
	isDebug = true
)

// CoinChaincode 实际操作类
type CoinChaincode struct {
}

type requestPacket struct {
	req    *Request
	caller *userItem
}

//Init initlaize chaincode
func (t *CoinChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("init")
	getBankManger(stub)
	getUserManger(stub)
	return shim.Success([]byte("init ok"))
}

//Invoke function call interface
func (t *CoinChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function != "callapi" {
		return shim.Error("chaincode function not equal :callapi")
	}
	if len(args) != 2 {
		return shim.Error("chaincode args count not equal 2")
	}
	req := &Request{}
	err := json.Unmarshal([]byte(args[0]), req)
	if err != nil {
		return shim.Error(err.Error())
	}
	umgr := getUserManger(stub)
	bmgr := getBankManger(stub)
	key := req.FromID
	packet := &requestPacket{
		req:    req,
		caller: nil,
	}
	if len(req.FromID) == 32 {
		user, err := umgr.getUser(key, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = t.checkSignature(args[0], args[1], user.PublicKey, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		packet.caller = user
	} else {
		err = t.checkSignature(args[0], args[1], req.FromID, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		packet.caller, _ = umgr.getUser(req.FromID, stub)
	}
	if packet.caller == nil && req.function != "adduser" {
		return shim.Error("FromID not exist")
	}
	return t.invoke(stub, packet, umgr, bmgr)
}

func (t *CoinChaincode) checkSignature(args string, sig string, key string, stub shim.ChaincodeStubInterface) error {
	if isDebug {
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

func (t *CoinChaincode) invoke(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	fmt.Println("invoke ", req.req.function, req.req.args)
	switch req.req.function {
	case "addbank":
		return t.addBank(stub, req, umgr, bmgr)
	case "adjustbanklimit":
		return t.adjustLimit(stub, req, umgr, bmgr)
	case "getbank":
		return t.getBank(stub, req, umgr, bmgr)
	case "setexchanemap":
		return t.setExchaneMap(stub, req, umgr, bmgr)
	case "adduser":
		return t.addUser(stub, req, umgr, bmgr)
	case "getuser":
		return t.getUser(stub, req, umgr, bmgr)
	case "cashin":
		return t.cashIn(stub, req, umgr, bmgr)
	case "cashout":
		return t.cashout(stub, req, umgr, bmgr)
	case "transfer":
		return t.transfer(stub, req, umgr, bmgr)
	case "exchange":
		return t.exchange(stub, req, umgr, bmgr)
	default:
		return shim.Error("invalid function:" + req.req.function)
	}
}

func (t *CoinChaincode) addUser(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg:[publickey]
	if len(req.req.args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	user, err := umgr.addUser(req.req.args[0], stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(user.toBuffer())
}

func (t *CoinChaincode) getUser(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg:[username]
	if len(req.req.args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	user, err := umgr.getUser(req.req.args[0], stub)
	if err != nil {
		return shim.Error("user name not found:" + req.req.args[0])
	}
	return shim.Success(user.toBuffer())
}

func (t *CoinChaincode) addBank(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg:[bankname,currency,chip,manger]
	args := req.req.args
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	if req.caller.UserType != userTypeRoot {
		return shim.Error("You can't call this function,access denyed")
	}
	user, err := umgr.getUser(args[3], stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if user.UserType != userTypeNormal {
		return shim.Error("the manger is used")
	}
	item := bankItem{
		BankName:   args[0],
		Currency:   args[1],
		Chip:       args[2],
		MangerName: user.ID,
	}
	bank, err := bmgr.addBank(stub, item)
	if err != nil {
		return shim.Error(err.Error())
	}
	_, err = umgr.upgradeUser(user.ID, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(bank.toBuffer())
}

func (t *CoinChaincode) adjustLimit(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg:[bankname,newvalue]
	args := req.req.args
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	value, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	if req.caller.UserType != userTypeRoot {
		return shim.Error("You can't call this function,access denyed")
	}
	bank, err := bmgr.lookupBankByName(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	bank.ChipLimit = value
	err = bmgr.save(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(bank.toBuffer())
}

func (t *CoinChaincode) getBank(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg: [bankname]
	args := req.req.args
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	bank, err := bmgr.lookupBankByName(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(bank.toBuffer())
}

func (t *CoinChaincode) setExchaneMap(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg:[isfixed,exchanemapjson]
	args := req.req.args
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	nVal := map[string]float64{}
	err := json.Unmarshal([]byte(args[1]), &nVal)
	if err != nil {
		return shim.Error(err.Error())
	}
	if req.caller.UserType != userTypeBankManger {
		return shim.Error("You can't call this function,access denyed")
	}
	bank, err := bmgr.lookupBankByMangerName(req.caller.ID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if args[0] == "true" {
		bank.FixedExchangeMap = nVal
	} else {
		bank.ExchangeMap = nVal
	}
	err = bmgr.save(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(bank.toBuffer())
}

func (t *CoinChaincode) cashIn(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg: [username,currency,amount]
	args := req.req.args
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	if req.caller.UserType != userTypeBankManger {
		return shim.Error("You are not a manger,can't use this function")
	}
	value, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	bank, err := bmgr.lookupBankByCurrency(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	if bank.MangerName != req.caller.ID {
		return shim.Error("You are not a manger of this bank " + bank.BankName)
	}
	user, err := umgr.getUser(args[0], stub)
	if user == nil {
		return shim.Error("user name not found:" + args[0])
	}
	if user.UserType != userTypeNormal {
		return shim.Error("You can only call this function on a normal user")
	}
	err = user.increaseBalance(args[1], value, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	bank.CurrencyCount += value
	err = bmgr.save(stub)
	if err != nil {
		user.decreaseBalance(args[1], value, stub)
		return shim.Error(err.Error())
	}
	return shim.Success(user.toBuffer())
}

func (t *CoinChaincode) cashout(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg: [username,currency,amount]
	args := req.req.args
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	if req.caller.UserType != userTypeBankManger {
		return shim.Error("You are not a manger,can't use this function")
	}
	value, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	user, err := umgr.getUser(args[0], stub)
	if user == nil {
		return shim.Error("user name not found:" + args[0])
	}
	bank, err := bmgr.lookupBankByCurrency(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	if bank.MangerName != req.caller.ID {
		return shim.Error("You are not a manger of this bank " + bank.BankName)
	}
	err = user.decreaseBalance(bank.Currency, value, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	bank.CurrencyCount -= value
	err = bmgr.save(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(user.toBuffer())
}

func (t *CoinChaincode) exchange(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg: [fromcurrency,tocurrency,amount,isfixedrate]
	args := req.req.args
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	if req.caller.UserType != userTypeNormal {
		return shim.Error("You are manger,can't use this function")
	}
	value, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	fromToken := false
	toToken := false
	frombank, err := bmgr.lookupBankByCurrency(args[0])
	if err != nil {
		frombank, err = bmgr.lookupBankByChip(args[0])
		if err != nil {
			return shim.Error(args[0] + " not fond")
		}
		fromToken = true
	}

	if fromToken && args[1] == frombank.Currency {
		//chipA->currencyA
		err = req.caller.decreaseBalance(args[0], value, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		req.caller.increaseBalance(args[1], value, stub)
		frombank.UsedChip -= value
		frombank.CurrencyCount += value
		bmgr.save(stub)
		return shim.Success(req.caller.toBuffer())
	}
	if !fromToken && args[1] == frombank.Chip {
		//currencyA->chipA
		if frombank.UsedChip+value > frombank.ChipLimit {
			return shim.Error(frombank.BankName + " bank out of limit")
		}
		err = req.caller.decreaseBalance(args[0], value, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		req.caller.increaseBalance(args[1], value, stub)
		frombank.UsedChip += value
		frombank.CurrencyCount -= value
		bmgr.save(stub)
		return shim.Success(req.caller.toBuffer())
	}
	tobank, err := bmgr.lookupBankByCurrency(args[1])
	if err != nil {
		tobank, err = bmgr.lookupBankByChip(args[1])
		if err != nil {
			return shim.Error(args[1] + " not fond")
		}
		toToken = true
	}
	isFixed := false
	if fromToken && toToken {
		isFixed = true
	}
	if args[3] == "true" {
		isFixed = true
	}
	rate := float64(0.0)
	if isFixed {
		rate = tobank.FixedExchangeMap[frombank.Currency+"2"+tobank.Currency]
	} else {
		rate = tobank.ExchangeMap[frombank.Currency+"2"+tobank.Currency]
	}
	if rate == 0.0 {
		return shim.Error("exchange rate not set for " + frombank.Currency + "2" + tobank.Currency)
	}
	err = req.caller.decreaseBalance(args[0], value, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	toValue := int((float64(value) * rate))
	if toToken {
		if tobank.UsedChip+toValue < tobank.ChipLimit {
			req.caller.increaseBalance(args[0], value, stub)
			return shim.Error(tobank.BankName + " out of limit")
		}
		tobank.UsedChip += toValue
	} else {
		tobank.CurrencyCount += toValue
	}
	if fromToken {
		frombank.UsedChip -= value
	} else {
		frombank.CurrencyCount -= value
	}
	if isFixed {
		req.caller.increaseLockedBalance(args[1], toValue, stub)
	} else {
		req.caller.increaseBalance(args[1], toValue, stub)
	}
	bmgr.save(stub)
	return shim.Success(req.caller.toBuffer())
}

func (t *CoinChaincode) transfer(stub shim.ChaincodeStubInterface, req *requestPacket, umgr *userManger, bmgr *bankManger) pb.Response {
	//arg: [touser,currency,amount,islocked]
	args := req.req.args
	fromUser := req.caller
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	value, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	toUser, err := umgr.getUser(args[0], stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if fromUser.UserType == userTypeNormal && toUser.UserType == userTypeNormal {
		//p2psend
		if args[3] == "true" {
			err = fromUser.decreaseLockedBalance(args[1], value, stub)
			if err != nil {
				return shim.Error(err.Error())
			}
			toUser.increaseLockedBalance(args[1], value, stub)
			return shim.Success(fromUser.toBuffer())
		}
		err = fromUser.decreaseBalance(args[1], value, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		toUser.increaseBalance(args[1], value, stub)
		return shim.Success(fromUser.toBuffer())
	}
	bank, err := bmgr.lookupBankByChip(args[1])
	if err != nil {
		return shim.Error("transfer currency on this user type not allowed")
	}

	if fromUser.UserType == userTypeBankManger && toUser.UserType == userTypeNormal {
		//chipreceive
		if bank.UsedChip+value > bank.ChipLimit {
			return shim.Error(bank.BankName + " chip out of limit")
		}
		bank.UsedChip += value
		toUser.increaseBalance(args[1], value, stub)
		bmgr.save(stub)
		return shim.Success(toUser.toBuffer())
	}
	if fromUser.UserType == userTypeNormal && toUser.UserType == userTypeBankManger {
		//chippay
		err = fromUser.decreaseBalance(args[1], value, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		bank.UsedChip -= value
		bmgr.save(stub)
		return shim.Success(fromUser.toBuffer())
	}
	return shim.Error("transfer from not allow for user type")
}

func main() {
	err := shim.Start(new(CoinChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

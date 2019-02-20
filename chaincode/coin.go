package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

const (
	DEBUG = true
)

type CoinChaincode struct {
}

func (t *CoinChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("init")
	GetBankManger(stub)
	GetUserManger(stub)
	return shim.Success([]byte("init ok"))
}

func (t *CoinChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoke")
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Invoke ", function, args[0])
	switch function {
	case "addbank":
		return t.addBank(stub, args)
	case "addbanklimit":
		return t.addBankAmount(stub, args)
	case "chippay":
		return t.chipPay(stub, args)
	case "getbank":
		return t.getBankInfomation(stub, args)
	case "adduser":
		return t.addUser(stub, args)
	case "getuser":
		return t.getUser(stub, args)
	case "cashin":
		return t.cashIn(stub, args)
	case "cashout":
		return t.cashout(stub, args)
	case "transfer":
		return t.cashout(stub, args)
	case "issue":
		return t.issue(stub, args)
	default:
		return shim.Error("invalid function:" + function)
	}
}

func (t *CoinChaincode) addBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg:[bankname,currency,chip,exchanger]
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	mgr := GetBankManger(stub)
	item := BankItem{
		BankName:  args[0],
		Currency:  args[1],
		Chip:      args[2],
		Exchanger: args[3],
	}
	bank, err := mgr.addBank(stub, item)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer mgr.save(stub)
	umgr := GetUserManger(stub)

	_, err = umgr.addUser(item.Exchanger, item.BankName, 2, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer umgr.save(stub)
	return shim.Success(bank.ToBuffer())
}

func (t *CoinChaincode) addBankAmount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg:[bankname,addvalue]
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	value, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	mgr := GetBankManger(stub)
	bank, err := mgr.lookupBankByName(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	bank.TotalAmount += value
	mgr.save(stub)
	return shim.Success(bank.ToBuffer())
}

func (t *CoinChaincode) getBankInfomation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg: [bankname]
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	mgr := GetBankManger(stub)
	bank, err := mgr.lookupBankByName(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(bank.ToBuffer())
}

func (t *CoinChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg:[username,bankname,type]
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	mgr := GetBankManger(stub)
	_, err := mgr.lookupBankByName(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	umgr := GetUserManger(stub)
	uType := 1
	if args[2] == "2" {
		uType = 2
	}
	item, err := umgr.addUser(args[0], args[1], uType, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	umgr.save(stub)
	return shim.Success(item.ToBuffer())
}

func (t *CoinChaincode) getUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg:[username]
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	umgr := GetUserManger(stub)
	user := umgr.getUser(args[0], stub)
	if user == nil {
		return shim.Error("user name not found:" + args[0])
	}
	return shim.Success(user.ToBuffer())
}

func (t *CoinChaincode) cashIn(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg: [username,currency,amount]
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	value, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	umgr := GetUserManger(stub)
	user := umgr.getUser(args[0], stub)
	if user == nil {
		return shim.Error("user name not found:" + args[0])
	}
	mgr := GetBankManger(stub)
	_, err = mgr.lookupBankByCurrency(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	err = user.increaseBalance(args[1], value, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(user.ToBuffer())
}

func (t *CoinChaincode) chipPay(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg: [bankname,username,currency,amount,islocked]
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	value, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error(err.Error())
	}
	umgr := GetUserManger(stub)
	user := umgr.getUser(args[1], stub)
	if user == nil {
		return shim.Error("user name not found:" + args[1])
	}
	mgr := GetBankManger(stub)
	bank, err := mgr.lookupBankByName(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if bank.Chip != args[2] {
		return shim.Error("You can ony use :" + bank.Chip + " on " + bank.BankName)
	}
	if args[4] == "true" {
		err = user.decreaseLockedBalance(bank.Chip, value, stub)
		if err != nil { //透支
			return shim.Error(err.Error())
		}
	} else {
		err = user.decreaseBalance(bank.Chip, value, stub)
		if err != nil { //透支
			return shim.Error(err.Error())
		}
	}
	bank.UsedAmount -= value
	bank.TotalAmount += value
	mgr.save(stub)
	return shim.Success(user.ToBuffer())
}

func (t *CoinChaincode) issue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg: [bankname,username,currency,amount]
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	value, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error(err.Error())
	}
	umgr := GetUserManger(stub)
	user := umgr.getUser(args[1], stub)
	if user == nil {
		return shim.Error("user name not found:" + args[1])
	}
	mgr := GetBankManger(stub)
	bank, err := mgr.lookupBankByName(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if bank.Currency != args[2] {
		return shim.Error("bankname and currency not match")
	}
	//额度不够
	if bank.TotalAmount < value {
		return shim.Error("bank limit overflow")
	}
	err = user.decreaseBalance(bank.Currency, value, stub)
	if err != nil { //透支
		return shim.Error(err.Error())
	}
	err = user.increaseBalance(bank.Chip, value, stub)
	if err != nil {
		user.increaseBalance(bank.Currency, value, stub)
		return shim.Error(err.Error())
	}
	bank.UsedAmount += value
	bank.TotalAmount -= value
	mgr.save(stub)
	return shim.Success(user.ToBuffer())
}

func (t *CoinChaincode) cashout(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg: [username,bankname,fromcurrency,dstcurrency,amount]
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	value, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error(err.Error())
	}
	umgr := GetUserManger(stub)
	user := umgr.getUser(args[0], stub)
	if user == nil {
		return shim.Error("user name not found:" + args[0])
	}
	mgr := GetBankManger(stub)
	bank, err := mgr.lookupBankByCurrency(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	if bank.Currency != args[2] {
		return shim.Error("bankname and currency not match")
	}
	if args[2] == args[3] { //同currency 提现，
		err = user.decreaseBalance(bank.Currency, value, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(user.ToBuffer())
	} else { //用USD 提取HKD？
		return shim.Error("cashout " + args[2] + "==>" + args[3] + " not implement")
	}
	return shim.Success(user.ToBuffer())
}

func (t *CoinChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//arg: [fromuser,touser,currency,amount,islocked]
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	value, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error(err.Error())
	}
	umgr := GetUserManger(stub)
	fromuser := umgr.getUser(args[0], stub)
	if fromuser == nil {
		return shim.Error("user name not found:" + args[0])
	}
	touser := umgr.getUser(args[1], stub)
	if touser == nil {
		return shim.Error("user name not found:" + args[1])
	}
	if args[4] == "true" {
		err = fromuser.decreaseLockedBalance(args[2], value, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = touser.increaseLockedBalance(args[2], value, stub)
		if err != nil {
			fromuser.increaseLockedBalance(args[2], value, stub)
			return shim.Error(err.Error())
		}
	} else {
		err = fromuser.decreaseBalance(args[2], value, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = touser.increaseBalance(args[2], value, stub)
		if err != nil {
			fromuser.increaseBalance(args[2], value, stub)
			return shim.Error(err.Error())
		}
	}
	return shim.Success(fromuser.ToBuffer())
}

func main() {
	err := shim.Start(new(CoinChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

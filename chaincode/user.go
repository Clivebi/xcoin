package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	keyUserManger      = "user_manger_root"
	userTypeRoot       = 0
	userTypeBankManger = 1
	userTypeNormal     = 2
)

//UserItem 存储用户信息结构
type userItem struct {
	PublicKey     string         `json:"pub_key"`
	ID            string         `json:"id"`
	UserType      int            `json:"type"`
	Balance       map[string]int `json:"balance"`
	LockedBalance map[string]int `json:"lockedbalance"`
}

// ToBuffer 序列化
func (o userItem) toBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

func (o userItem) save(stub shim.ChaincodeStubInterface) error {
	key := "members_" + o.ID
	return stub.PutState(key, o.toBuffer())
}

func (o *userItem) increaseBalance(sym string, value int, stub shim.ChaincodeStubInterface) error {
	old := o.Balance[sym]
	old += value
	o.Balance[sym] = old
	o.save(stub)
	return nil
}

func (o *userItem) decreaseBalance(sym string, value int, stub shim.ChaincodeStubInterface) error {
	old := o.Balance[sym]
	if old < value {
		return errors.New("out of balance")
	}
	old -= value
	o.Balance[sym] = old
	o.save(stub)
	return nil
}

func (o *userItem) increaseLockedBalance(sym string, value int, stub shim.ChaincodeStubInterface) error {
	old := o.LockedBalance[sym]
	old += value
	o.LockedBalance[sym] = old
	o.save(stub)
	return nil
}

func (o *userItem) decreaseLockedBalance(sym string, value int, stub shim.ChaincodeStubInterface) error {
	old := o.LockedBalance[sym]
	if old < value {
		return errors.New("out of balance")
	}
	old -= value
	o.LockedBalance[sym] = old
	o.save(stub)
	return nil
}

type userManger struct {
	Root        string   `json:"root"`
	BankMangers []string `json:"bankmangers"` //所有管理员
	Users       []string `json:"users"`       //所有用户
}

func (o userManger) toBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

func getUserManger(stub shim.ChaincodeStubInterface) *userManger {
	obj := &userManger{
		Users:       []string{},
		BankMangers: []string{},
	}
	buf, _ := stub.GetState(keyUserManger)
	if buf != nil {
		json.Unmarshal(buf, obj)
	} else {
		obj.save(stub)
	}
	return obj
}

func (o *userManger) save(stub shim.ChaincodeStubInterface) error {
	return stub.PutState(keyUserManger, o.toBuffer())
}

func (o *userManger) getUserFromID(name string, stub shim.ChaincodeStubInterface) (*userItem, error) {
	key := "members_" + name
	value, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	obj := &userItem{}
	err = json.Unmarshal(value, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (o *userManger) getUserID(publickey string) string {
	return EncodeWalletAddress(Base58Decode(publickey))
}

func (o *userManger) getUser(key string, stub shim.ChaincodeStubInterface) (*userItem, error) {
	if !IsWalletAddress(key) {
		key = o.getUserID(key)
	}
	return o.getUserFromID(key, stub)
}

func (o *userManger) addUser(publickey string, stub shim.ChaincodeStubInterface) (*userItem, error) {
	user, err := o.getUser(publickey, stub)
	if err == nil {
		return nil, errors.New("user alerady exist")
	}
	user = &userItem{
		PublicKey:     publickey,
		ID:            o.getUserID(publickey),
		UserType:      userTypeNormal,
		Balance:       map[string]int{},
		LockedBalance: map[string]int{},
	}
	if len(o.Root) == 0 {
		user.UserType = userTypeRoot
		o.Root = user.ID
	}
	o.Users = append(o.Users, user.ID)
	err = user.save(stub)
	if err != nil {
		return nil, err
	}
	err = o.save(stub)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (o *userManger) upgradeUser(key string, stub shim.ChaincodeStubInterface) (*userItem, error) {
	user, err := o.getUser(key, stub)
	if err != nil {
		return nil, err
	}
	if user.UserType != userTypeNormal {
		return user, nil
	}
	o.BankMangers = append(o.BankMangers, user.ID)
	user.UserType = userTypeBankManger
	if err := user.save(stub); err != nil {
		return nil, err
	}
	if err := o.save(stub); err != nil {
		return nil, err
	}
	return user, nil
}

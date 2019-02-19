package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	KEY_USER_MANGER = "user_manger_root"
)

type UserItem struct {
	UserName      string         `json:"name"`
	BirthBankName string         `json:"bank"`
	Type          int            `json:"type"`
	Balance       map[string]int `json:"balance"`
	LockedBalance map[string]int `json:"lockedbalance"`
}

func (o UserItem) ToBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

func (o UserItem) IsSeller() bool {
	return o.Type == 1
}

func (o UserItem) save(stub shim.ChaincodeStubInterface) {
	key := "members_" + o.UserName
	stub.PutState(key, o.ToBuffer())
}

func (o *UserItem) increaseBalance(sym string, value int, stub shim.ChaincodeStubInterface) error {
	old := o.Balance[sym]
	old += value
	o.Balance[sym] = old
	o.save(stub)
	return nil
}

func (o *UserItem) decreaseBalance(sym string, value int, stub shim.ChaincodeStubInterface) error {
	old := o.Balance[sym]
	if old < value {
		return errors.New("out of balance")
	}
	old -= value
	o.Balance[sym] = old
	o.save(stub)
	return nil
}

func (o *UserItem) increaseLockedBalance(sym string, value int, stub shim.ChaincodeStubInterface) error {
	old := o.LockedBalance[sym]
	old += value
	o.LockedBalance[sym] = old
	o.save(stub)
	return nil
}

func (o *UserItem) decreaseLockedBalance(sym string, value int, stub shim.ChaincodeStubInterface) error {
	old := o.LockedBalance[sym]
	if old < value {
		return errors.New("out of balance")
	}
	old -= value
	o.LockedBalance[sym] = old
	o.save(stub)
	return nil
}

type UserManger struct {
	Sellers  []string `json:"sellers"`   //所有管理员
	Users    []string `json:"users"`     //所有用户
	LogIndex int      `json:"log_index"` //日志索引，类似于mysql的自增长ID
}

func (o UserManger) ToBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

func GetUserManger(stub shim.ChaincodeStubInterface) *UserManger {
	obj := &UserManger{
		Users:    []string{},
		Sellers:  []string{},
		LogIndex: 0,
	}
	buf, _ := stub.GetState(KEY_USER_MANGER)
	if buf != nil {
		json.Unmarshal(buf, obj)
	} else {
		obj.save(stub)
	}
	return obj
}

func (o *UserManger) save(stub shim.ChaincodeStubInterface) {
	stub.PutState(KEY_USER_MANGER, o.ToBuffer())
}

func (o *UserManger) getUser(name string, stub shim.ChaincodeStubInterface) *UserItem {
	key := "members_" + name
	value, err := stub.GetState(key)
	if err != nil {
		return nil
	}
	obj := &UserItem{}
	err = json.Unmarshal(value, obj)
	if err != nil {
		return nil
	}
	return obj
}

func (o *UserManger) addUser(name string, bankname string, utype int, stub shim.ChaincodeStubInterface) (*UserItem, error) {
	user := o.getUser(name, stub)
	if user != nil {
		return nil, errors.New("user alerady exist")
	}
	if utype != 1 && utype != 2 {
		return nil, errors.New("user type error")
	}
	user = &UserItem{
		UserName:      name,
		BirthBankName: bankname,
		Balance:       map[string]int{},
		LockedBalance: map[string]int{},
		Type:          utype,
	}
	if utype == 2 {
		o.Sellers = append(o.Sellers, name)
	} else {
		o.Users = append(o.Users, name)
	}
	user.save(stub)
	o.save(stub)
	return user, nil
}

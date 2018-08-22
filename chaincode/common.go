package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	MANGER_KEY = "mangers"
)

type UserItem struct {
	ID        string `json:"id"`     //用户ID，md5（pubkey)
	Type      int    `json:"type"`   //用户类型，0 root 1 管理员 2 顾客
	Coin      int    `json:"coin"`   //当前筹码
	PubKey    string `json:"pubkey"` //用户公钥
	CoinBase  int    `json:"base"`   //从root预借的筹码,只有Type ==1 才有这两个字段
	CoinLimit int    `json:"limit"`  //可以从root预借的最大筹码
}

func (o UserItem) ToBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

func (o UserItem) IsRoot() bool {
	return o.Type == 0
}

func (o UserItem) IsSeller() bool {
	return o.Type == 1
}

func (o UserItem) save(stub shim.ChaincodeStubInterface) {
	key := "user_" + o.ID
	stub.PutState(key, o.ToBuffer())
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

func NewManger(stub shim.ChaincodeStubInterface) *UserManger {
	obj := &UserManger{
		Users:    []string{},
		Sellers:  []string{},
		LogIndex: 0,
	}
	buf, _ := stub.GetState(MANGER_KEY)
	if buf != nil {
		json.Unmarshal(buf, obj)
	} else {
		obj.save(stub)
	}
	return obj
}

func (o *UserManger) save(stub shim.ChaincodeStubInterface) {
	stub.PutState(MANGER_KEY, o.ToBuffer())
}

func (o *UserManger) getUser(ID string, stub shim.ChaincodeStubInterface) *UserItem {
	key := "user_" + ID
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

func (o *UserManger) publicKeyToUserID(pubkey string) string {
	hs := md5.New()
	hs.Write([]byte(pubkey))
	return hex.EncodeToString(hs.Sum(nil))
}

func (o *UserManger) AddUser(pubkey string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	ID := o.publicKeyToUserID(pubkey)
	user := o.getUser(ID, stub)
	if user != nil {
		return nil, errors.New("user alerady exist")
	}
	user = &UserItem{
		ID:        ID,
		PubKey:    pubkey,
		Coin:      0,
		Type:      2,
		CoinBase:  0,
		CoinLimit: 0,
	}
	if len(o.Sellers) == 0 {
		user.Type = 0
		o.Sellers = append(o.Sellers, ID)
	}
	o.Users = append(o.Users, ID)
	user.save(stub)
	o.save(stub)
	return user.ToBuffer(), nil
}

func (o *UserManger) GetUser(ID string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	user := o.getUser(ID, stub)
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user.ToBuffer(), nil
}

func (o *UserManger) UpgradeUser(Caller *UserItem, ID string, Limit int, stub shim.ChaincodeStubInterface) error {
	if !Caller.IsRoot() {
		return errors.New("aceess denyed")
	}
	user := o.getUser(ID, stub)
	if user == nil {
		return errors.New("user not found")
	}
	if user.Coin != 0 {
		return errors.New("user not empty")
	}
	if user.Type != 2 {
		return errors.New("user not empty")
	}
	user.CoinBase = 0
	user.Coin = 0
	user.Type = 1
	user.CoinLimit = Limit
	user.save(stub)
	o.Sellers = append(o.Sellers, user.ID)
	o.save(stub)
	return nil
}

func (o *UserManger) SetSellerLimit(Caller *UserItem, ID string, Limit int, stub shim.ChaincodeStubInterface) error {
	if !Caller.IsRoot() {
		return errors.New("aceess denyed")
	}
	user := o.getUser(ID, stub)
	if user == nil {
		return errors.New("user not found")
	}
	if !user.IsSeller() {
		return errors.New("user not a seller")
	}
	return errors.New("not implement")
}

func (o *UserManger) Send(Caller *UserItem, ID string, coin int, stub shim.ChaincodeStubInterface) error {
	user := o.getUser(ID, stub)
	if user == nil {
		return errors.New("user not found")
	}
	if Caller.IsRoot() {
		return errors.New("root not have any coin")
	}
	if coin < 0 {
		return errors.New("invliad coin value")
	}
	if Caller.Coin < coin {
		if !Caller.IsSeller() {
			return errors.New("out of coin")
		}
		if Caller.CoinBase+coin < Caller.CoinLimit {
			Caller.CoinBase += coin
			Caller.Coin += coin
		}
		if coin > Caller.CoinBase {
			return errors.New("seller arrive the limit coin")
		}
	}
	user.Coin += coin
	Caller.Coin -= coin
	user.save(stub)
	Caller.save(stub)
	return nil
}

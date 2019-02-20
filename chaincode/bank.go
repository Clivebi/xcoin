package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

const (
	KEY_BANK_MANGER = "bank_manger_root"
)

type BankItem struct {
	BankName    string `json:"bankname"`
	Currency    string `json:"currency"`
	Chip        string `json:"chip"`
	TotalAmount int    `json:"totalamount"`
	UsedAmount  int    `json:"usedamount"`
	Exchanger   string `json:"exchanger"`
}

func (o BankItem) ToBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

type BankManger struct {
	Banks    []*BankItem
	UsedKeys string
}

func (o BankManger) ToBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

func GetBankManger(stub shim.ChaincodeStubInterface) *BankManger {
	obj := &BankManger{}
	buf, _ := stub.GetState(KEY_BANK_MANGER)
	if buf != nil {
		json.Unmarshal(buf, obj)
	} else {
		obj.save(stub)
	}
	return obj
}

func (o *BankManger) save(stub shim.ChaincodeStubInterface) error {
	return stub.PutState(KEY_BANK_MANGER, o.ToBuffer())
}

//以下两接口检查命名冲突和添加命名
func (o *BankManger) checkUsedKeys(bank BankItem) error {
	key := "|" + bank.BankName + "|"
	if strings.Contains(o.UsedKeys, key) {
		return errors.New("bank name used:" + bank.BankName)
	}
	key = "|" + bank.Chip + "|"
	if strings.Contains(o.UsedKeys, key) {
		return errors.New("chip name used:" + bank.Chip)
	}
	key = "|" + bank.Currency + "|"
	if strings.Contains(o.UsedKeys, key) {
		return errors.New("currency used:" + bank.Currency)
	}
	key = "|" + bank.Exchanger + "|"
	if strings.Contains(o.UsedKeys, key) {
		return errors.New("exchanger name used:" + bank.Exchanger)
	}
	return nil
}

func (o *BankManger) addUsedKeys(bank BankItem) {
	key := "|" + bank.BankName + "|"
	key += bank.Chip
	key += "|"
	key += bank.Currency
	key += "|"
	key += bank.Exchanger
	key += "|"
	o.UsedKeys += key
}

func (o *BankManger) lookupBankByCurrency(currency string) (*BankItem, error) {
	for _, it := range o.Banks {
		if it.Currency == currency {
			return it, nil
		}
	}
	return nil, errors.New("bank not found")
}

func (o *BankManger) lookupBankByName(name string) (*BankItem, error) {
	for _, it := range o.Banks {
		if it.BankName == name {
			return it, nil
		}
	}
	return nil, errors.New("bank not found")
}

func (o *BankManger) lookupBankByChip(chip string) (*BankItem, error) {
	for _, it := range o.Banks {
		if it.Chip == chip {
			return it, nil
		}
	}
	return nil, errors.New("bank not found")
}

func (o *BankManger) lookupBankByExchanger(name string) (*BankItem, error) {
	for _, it := range o.Banks {
		if it.Exchanger == name {
			return it, nil
		}
	}
	return nil, errors.New("bank not found")
}

func (o *BankManger) addBank(stub shim.ChaincodeStubInterface, item BankItem) error {
	if err := o.checkUsedKeys(item); err != nil {
		return err
	}
	nit := &BankItem{
		BankName:    item.BankName,
		Currency:    item.Currency,
		Chip:        item.Chip,
		TotalAmount: item.TotalAmount,
		UsedAmount:  0,
		Exchanger:   item.Exchanger,
	}
	o.Banks = append(o.Banks, nit)
	o.addUsedKeys(item)
	return o.save(stub)
}

package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	KEY_BANK_MANGER = "bank_manger_root"
)

type BankItem struct {
	BankName    string `json:"bankname"`
	Currency    string `json:"currenty"`
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
	Banks []*BankItem
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
	_, err := o.lookupBankByCurrency(item.Currency)
	if err == nil {
		return errors.New("Currency exist")
	}
	_, err = o.lookupBankByName(item.BankName)
	if err == nil {
		return errors.New("bank name exist")
	}
	_, err = o.lookupBankByChip(item.Chip)
	if err == nil {
		return errors.New("chip name exist")
	}
	_, err = o.lookupBankByExchanger(item.Exchanger)
	if err == nil {
		return errors.New("exchanger name exist")
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
	return o.save(stub)
}

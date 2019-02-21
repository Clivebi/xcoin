package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

const (
	keyBankManger = "bank_manger_root"
)

type bankItem struct {
	BankName         string             `json:"bankname"`
	Currency         string             `json:"currency"`
	Chip             string             `json:"chip"`
	ChipLimit        int                `json:"chiplimit"`
	UsedChip         int                `json:"chipused"`
	CurrencyCount    int                `json:"currencyCount"`
	MangerName       string             `json:"mangername"`
	ExchangeMap      map[string]float64 `json:"exchangemap"`
	FixedExchangeMap map[string]float64 `json:"fiexedexchangemap"`
}

func (o bankItem) toBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

type bankManger struct {
	Banks    []*bankItem
	UsedKeys string
}

func (o bankManger) toBuffer() []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return buf
}

func getBankManger(stub shim.ChaincodeStubInterface) *bankManger {
	obj := &bankManger{}
	buf, _ := stub.GetState(keyBankManger)
	if buf != nil {
		json.Unmarshal(buf, obj)
	} else {
		obj.save(stub)
	}
	return obj
}

func (o *bankManger) save(stub shim.ChaincodeStubInterface) error {
	return stub.PutState(keyBankManger, o.toBuffer())
}

//以下两接口检查命名冲突和添加命名
func (o *bankManger) checkUsedKeys(bank bankItem) error {
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
	return nil
}

func (o *bankManger) addUsedKeys(bank bankItem) {
	key := "|" + bank.BankName + "|"
	key += bank.Chip
	key += "|"
	key += bank.Currency
	key += "|"
	o.UsedKeys += key
}

func (o *bankManger) lookupBankByCurrency(currency string) (*bankItem, error) {
	for _, it := range o.Banks {
		if it.Currency == currency {
			return it, nil
		}
	}
	return nil, errors.New("bank not found")
}

func (o *bankManger) lookupBankByName(name string) (*bankItem, error) {
	for _, it := range o.Banks {
		if it.BankName == name {
			return it, nil
		}
	}
	return nil, errors.New("bank not found")
}

func (o *bankManger) lookupBankByChip(chip string) (*bankItem, error) {
	for _, it := range o.Banks {
		if it.Chip == chip {
			return it, nil
		}
	}
	return nil, errors.New("bank not found")
}

func (o *bankManger) lookupBankByMangerName(name string) (*bankItem, error) {
	for _, it := range o.Banks {
		if it.MangerName == name {
			return it, nil
		}
	}
	return nil, errors.New("bank not found")
}

func (o *bankManger) addBank(stub shim.ChaincodeStubInterface, item bankItem) (*bankItem, error) {
	if err := o.checkUsedKeys(item); err != nil {
		return nil, err
	}
	nit := &bankItem{
		BankName:      item.BankName,
		Currency:      item.Currency,
		Chip:          item.Chip,
		ChipLimit:     0,
		UsedChip:      0,
		CurrencyCount: 0,
		MangerName:    item.MangerName,
	}
	o.Banks = append(o.Banks, nit)
	o.addUsedKeys(item)
	return nit, o.save(stub)
}

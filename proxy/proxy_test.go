package proxy

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var publickkeys []crypto.PublicKey
var privatekeys []crypto.PrivateKey
var userids []string

const (
	rootUser      = 0
	mangerOfBankA = 1
	mangerOfBankB = 2
	normalUserA   = 3
	normalUserB   = 4
	apiURL        = "http://127.0.0.1:8789/callapi.do"
)

type Request struct {
	Time     int64    `json:"timestamp"` //时间戳
	FromID   string   `json:"fromid"`    //调用者，用户的ID或者public key
	Function string   `json:"func"`      //调用函数
	Args     []string `json:"args"`      //调用参数
}

type sendBuffer struct {
	Req string `json:"req"`
	Sig string `json:"sig"`
}

func initEnv() {
	publickkeys = make([]crypto.PublicKey, 5)
	privatekeys = make([]crypto.PrivateKey, 5)
	userids = make([]string, 5)

	for i := 0; i < 5; i++ {
		key, _ := rsa.GenerateKey(rand.Reader, 2048)
		publickkeys[i] = key.Public()
		privatekeys[i] = key
	}
}

func publickeyText(role int) string {
	buf := x509.MarshalPKCS1PublicKey(publickkeys[role].(*rsa.PublicKey))
	return base64.StdEncoding.EncodeToString(buf)
}

func signText(pk crypto.PrivateKey, text string) (string, error) {
	hash := sha256.Sum256([]byte(text))
	buf, err := rsa.SignPKCS1v15(rand.Reader, pk.(*rsa.PrivateKey), crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

//GetPublickeyText 从私钥获取文本格式的公钥
func GetPublickeyText(pk crypto.PrivateKey) string {
	buf := x509.MarshalPKCS1PublicKey(pk.(*rsa.PrivateKey).Public().(*rsa.PublicKey))
	return base64.StdEncoding.EncodeToString(buf)
}

//CallAPI 调用合约功能
//pk 发起函数调用的调用者的私钥
//function args 将要调用合约的功能和参数
func CallAPI(pk crypto.PrivateKey, function string, args []string, t *testing.T) (string, error) {
	req := &Request{
		FromID:   GetPublickeyText(pk),
		Time:     time.Now().Unix(),
		Function: function,
		Args:     args,
	}
	buf, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	sbuf := &sendBuffer{
		Req: string(buf),
		Sig: "",
	}
	sbuf.Sig, err = signText(pk, sbuf.Req)
	if err != nil {
		return "", err
	}
	buf, _ = json.Marshal(sbuf)
	rd := bytes.NewReader(buf)
	rsp, err := http.Post(apiURL, "application/json", rd)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != 200 {
		return "", errors.New("HTTP ERROR :" + rsp.Status)
	}
	t.Log(sbuf.Req)
	buf, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return string(buf), err
}

func callAPI(role int, function string, args []string, t *testing.T) (string, error) {
	text, err := CallAPI(privatekeys[role], function, args, t)
	time.Sleep(time.Second * 5)
	return text, err
}

func TestUser(t *testing.T) {
	initEnv()
	t.Log("add root")
	out, err := callAPI(rootUser, "adduser", []string{publickeyText(rootUser)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log("add bank manger A")
	out, err = callAPI(mangerOfBankA, "adduser", []string{publickeyText(mangerOfBankA)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log("add bank manger B")
	out, err = callAPI(mangerOfBankB, "adduser", []string{publickeyText(mangerOfBankB)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log("add normalUserA")
	out, err = callAPI(normalUserA, "adduser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log("add normalUserB")
	out, err = callAPI(normalUserB, "adduser", []string{publickeyText(normalUserB)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log("get normalUserA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserB)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
}

func TestBank(t *testing.T) {
	t.Log(" add bank must success")
	out, err := callAPI(rootUser, "addbank", []string{"bankA", "USD", "TokenA", publickeyText(mangerOfBankA)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log(" add bank must failed")
	out, err = callAPI(normalUserA, "addbank", []string{"bankB", "HKD", "TokenB", publickeyText(mangerOfBankB)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log(" add bank must success")
	out, err = callAPI(rootUser, "addbank", []string{"bankB", "HKD", "TokenB", publickeyText(mangerOfBankB)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)

	t.Log(" get mangerOfBankB info")
	out, err = callAPI(mangerOfBankB, "getuser", []string{publickeyText(mangerOfBankB)}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)

	t.Log(" adjust bankA limit")
	out, err = callAPI(rootUser, "adjustbanklimit", []string{"bankA", "1000000"}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)

	t.Log(" adjust bankB limit")
	out, err = callAPI(rootUser, "adjustbanklimit", []string{"bankB", "1000000"}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)

	t.Log(" get bank A info")
	out, err = callAPI(mangerOfBankA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
}

func TestCashin(t *testing.T) {
	t.Log("cashin")
	out, err := callAPI(mangerOfBankA, "cashin", []string{publickeyText(normalUserA), "USD", "2000"}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log("after cashin query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after cashin query bankA")
	out, err = callAPI(normalUserA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}

func TestCashout(t *testing.T) {
	t.Log("cashout")
	out, err := callAPI(mangerOfBankA, "cashout", []string{publickeyText(normalUserA), "USD", "100"}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log("after cashout query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after cashout query bankA")
	out, err = callAPI(normalUserA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}

func TestP2PSend(t *testing.T) {
	t.Log("P2PSend")
	out, err := callAPI(normalUserA, "transfer", []string{publickeyText(normalUserB), "USD", "100", "false"}, t)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
	t.Log("after P2PSend query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after P2PSend query userB")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserB)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}

func TestIssue(t *testing.T) {
	t.Log("Issue")
	out, err := callAPI(normalUserA, "exchange", []string{"USD", "TokenA", "500", "false"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after Issue query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after Issue query bankA")
	out, err = callAPI(normalUserA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}
func TestChipPay(t *testing.T) {
	t.Log("ChipPay")
	out, err := callAPI(normalUserA, "transfer", []string{publickeyText(mangerOfBankA), "TokenA", "200", "false"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ChipPay query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ChipPay query bankA")
	out, err = callAPI(normalUserA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}

func TestChipReceive(t *testing.T) {
	t.Log("ChipReceive")
	out, err := callAPI(mangerOfBankA, "transfer", []string{publickeyText(normalUserA), "TokenA", "100", "false"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ChipReceive query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ChipReceive query bankA")
	out, err = callAPI(normalUserA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}

func TestSetExchangeMap(t *testing.T) {
	t.Log("SetExchangeMap")
	out, err := callAPI(mangerOfBankB, "setexchanemap", []string{"false", "{\"USD2HKD\":7.0,\"HKD2USD\":0.14}"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	out, err = callAPI(mangerOfBankB, "setexchanemap", []string{"true", "{\"USD2HKD\":8.0,\"HKD2USD\":0.12}"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}
func TestExchangeC2C(t *testing.T) {
	t.Log("ExchangeC2C")
	out, err := callAPI(normalUserA, "exchange", []string{"USD", "HKD", "100", "false"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeC2C query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeC2C query bankA")
	out, err = callAPI(normalUserA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeC2C query bankB")
	out, err = callAPI(normalUserA, "getbank", []string{"bankB"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}

func TestExchangeC2T(t *testing.T) {
	t.Log("ExchangeC2T")
	out, err := callAPI(normalUserA, "exchange", []string{"USD", "TokenB", "100", "false"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeC2T query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeC2T query bankA")
	out, err = callAPI(normalUserA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeC2T query bankB")
	out, err = callAPI(normalUserA, "getbank", []string{"bankB"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}

func TestExchangeT2T(t *testing.T) {
	t.Log("ExchangeT2T")
	out, err := callAPI(normalUserA, "exchange", []string{"TokenA", "TokenB", "100", "false"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeT2T query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeT2T query bankA")
	out, err = callAPI(normalUserA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeT2T query bankB")
	out, err = callAPI(normalUserA, "getbank", []string{"bankB"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}

func TestExchangeT2C(t *testing.T) {
	t.Log("ExchangeT2C")
	out, err := callAPI(normalUserA, "exchange", []string{"TokenA", "USD", "100", "false"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeT2C query userA")
	out, err = callAPI(normalUserA, "getuser", []string{publickeyText(normalUserA)}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeT2C query bankA")
	out, err = callAPI(normalUserA, "getbank", []string{"bankA"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Log("after ExchangeT2C query bankB")
	out, err = callAPI(normalUserA, "getbank", []string{"bankB"}, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}

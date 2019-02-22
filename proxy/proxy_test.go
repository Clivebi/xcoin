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

func signText(role int, text string) (string, error) {
	hash := sha256.Sum256([]byte(text))
	buf, err := rsa.SignPKCS1v15(rand.Reader, privatekeys[role].(*rsa.PrivateKey), crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

func callAPI(role int, function string, args []string, t *testing.T) (string, error) {
	req := &Request{
		FromID:   publickeyText(role),
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
	sbuf.Sig, err = signText(role, sbuf.Req)
	if err != nil {
		return "", err
	}
	buf, _ = json.Marshal(sbuf)
	rd := bytes.NewReader(buf)
	rsp, err := http.Post("http://127.0.0.1:8789/callapi.do", "application/json", rd)
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
	time.Sleep(5 * time.Second)
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

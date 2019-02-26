package proxy

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

//Request 函数调用结构体
type Request struct {
	Time     int64    `json:"timestamp"` //时间戳
	FromID   string   `json:"fromid"`    //调用者，用户的ID或者public key
	Function string   `json:"func"`      //调用函数
	Args     []string `json:"args"`      //调用参数
}

//Signature 签名信息
type Signature struct {
	Caller  string `json:"caller"`
	OptUser string `json:"optuser"`
}

type sendBuffer struct {
	Req string `json:"req"`
	Sig string `json:"sig"`
}

//PublicKeyToString get base64 encode publick key
func PublicKeyToString(publicKey *rsa.PublicKey) string {
	buf := x509.MarshalPKCS1PublicKey(publicKey)
	return Base58Encode(buf)
}

//PublicKeyToID convert publick key to user ID
func PublicKeyToID(publicKey *rsa.PublicKey) string {
	buf := x509.MarshalPKCS1PublicKey(publicKey)
	return EncodeWalletAddress(buf)
}

//NewRequest create request from caller
func NewRequest(callID string, function string, args []string) (string, error) {
	req := &Request{
		Time:     time.Now().Unix(),
		FromID:   callID,
		Function: function,
		Args:     args,
	}
	buf, err := json.Marshal(req)
	return string(buf), err
}

//SignRequest get request signature by privatekey
func SignRequest(request string, privatekey *rsa.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(request))
	buf, err := rsa.SignPKCS1v15(rand.Reader, privatekey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return Base58Encode(buf), nil
}

//CallAPI send request and receive response
// example
//   req,_ := NewRequest(PublicKeyToID(publicKey),"adduser",PublicKeyToString(publicKey))
//   callersig,_ :=SignRequest(req,privateKey)
//   rsp,err := CallAPI("http://127.0.0.1:8789/callapi.do",req,callersig,"")
func CallAPI(apiURI string, request string, callersig string, optusersig string) (string, error) {
	sig := &Signature{
		Caller:  callersig,
		OptUser: optusersig,
	}
	buf, err := json.Marshal(sig)
	if err != nil {
		return "", err
	}
	sbuf := &sendBuffer{
		Req: request,
		Sig: string(buf),
	}
	rbytes, err := json.Marshal(sbuf)
	if err != nil {
		return "", err
	}
	rd := bytes.NewReader(rbytes)
	rsp, err := http.Post(apiURI, "application/json", rd)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != 200 {
		return "", errors.New("HTTP ERROR " + rsp.Status)
	}
	buf, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return string(buf), err
}

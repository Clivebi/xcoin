package proxy

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//SignWithECDSA signature request with ecdsa
func SignWithECDSA(text []byte, pk *ecdsa.PrivateKey) (string, error) {
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(text)
	texthash := sha256_h.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, pk, texthash)
	if err != nil {
		return "", err
	}
	rt, err := r.MarshalText()
	if err != nil {
		return "", err
	}
	st, err := s.MarshalText()
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	_, err = w.Write([]byte(string(rt) + "+" + string(st)))
	if err != nil {
		return "", err
	}
	w.Flush()
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

//SignWithRSA signature request with rsa
func SignWithRSA(request string, privatekey *rsa.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(request))
	buf, err := rsa.SignPKCS1v15(rand.Reader, privatekey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

//StringToPublicKey convert string to rsa or ecdsa public key
func StringToPublicKey(text string) (interface{}, error) {
	buf, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKIXPublicKey(buf)
}

//PublicKeyToString convert rsa or ecdsa public key to string
func PublicKeyToString(pubk interface{}) (string, error) {
	buf, err := x509.MarshalPKIXPublicKey(pubk)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

//Request 函数调用结构体
type Request struct {
	Time     int64    `json:"timestamp"` //时间戳
	Function string   `json:"func"`      //调用函数
	Args     []string `json:"args"`      //调用参数
}

type sendBuffer struct {
	Req string `json:"req"`
	Sig string `json:"sig"`
}

//NewRequest create request from caller
func NewRequest(function string, args []string) (string, error) {
	req := &Request{
		Time:     time.Now().Unix(),
		Function: function,
		Args:     args,
	}
	buf, err := json.Marshal(req)
	return string(buf), err
}

//SignRequest get request signature by privatekey
func SignRequest(request string, privateKey interface{}) (string, error) {
	switch pk := privateKey.(type) {
	case *rsa.PrivateKey:
		return SignWithRSA(request, pk)
	case *ecdsa.PrivateKey:
		return SignWithECDSA([]byte(request), pk)
	default:
		return "", errors.New("unsupport private key ")
	}
}

//CallAPI send request and receive response
// example
//   req,_ := NewRequest(PublicKeyToID(publicKey),"adduser",PublicKeyToString(publicKey))
//   callersig,_ :=SignRequest(req,privateKey)
//   rsp,err := CallAPI("http://127.0.0.1:8789/callapi.do",req,callersig,"")
func CallAPI(apiURI string, request string, callersig string) (*Response, error) {
	sbuf := &sendBuffer{
		Req: request,
		Sig: callersig,
	}
	rbytes, err := json.Marshal(sbuf)
	if err != nil {
		return nil, err
	}
	rd := bytes.NewReader(rbytes)
	rsp, err := http.Post(apiURI, "application/json", rd)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != 200 {
		return nil, errors.New("HTTP ERROR " + rsp.Status)
	}
	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	robj := &Response{}
	err = json.Unmarshal(buf, robj)
	return robj, err
}

//RegisgerWallet register wallet to the chain
func RegisgerWallet(apiURL string, walletAddress string, publickey interface{}, group string) (*Wallet, error) {
	pkText, err := PublicKeyToString(publickey)
	if err != nil {
		return nil, err
	}
	request, err := NewRequest("register", []string{walletAddress, pkText, group})
	if err != nil {
		return nil, err
	}
	rsp, err := CallAPI(apiURL, request, "")
	if err != nil {
		return nil, err
	}
	if rsp.ErrorMessage != "success" {
		return nil, errors.New(rsp.ErrorMessage)
	}
	return rsp.Payload, nil
}

//GetWallet get wallet infomation from chain
func GetWallet(apiURL string, walletAddress string) (*Wallet, error) {
	request, err := NewRequest("getwallet", []string{walletAddress})
	if err != nil {
		return nil, err
	}
	rsp, err := CallAPI(apiURL, request, "")
	if err != nil {
		return nil, err
	}
	if rsp.ErrorMessage != "success" {
		return nil, errors.New(rsp.ErrorMessage)
	}
	return rsp.Payload, nil
}

//Send send token from fromwallet to toWallet
func Send(apiURL string, fromWallet string, toWallet string, amount float64, privateKeyOfFrom interface{}) (*Wallet, error) {
	request, err := NewRequest("send", []string{fromWallet, toWallet, strconv.FormatFloat(amount, 'g', 5, 32)})
	if err != nil {
		return nil, err
	}
	sig, err := SignRequest(request, privateKeyOfFrom)
	if err != nil {
		return nil, err
	}
	rsp, err := CallAPI(apiURL, request, sig)
	if err != nil {
		return nil, err
	}
	if rsp.ErrorMessage != "success" {
		return nil, errors.New(rsp.ErrorMessage)
	}
	return rsp.Payload, nil
}

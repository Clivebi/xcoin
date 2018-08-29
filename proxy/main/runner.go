package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/xcoin/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type CoinServer struct {
	runner *proxy.AppRunner
}

func (s *CoinServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		r.Form = make(url.Values)
		r.PostForm = make(url.Values)
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.NotFound(w, r)
			log.Println("read json body failed", err)
			return
		}
		if buf == nil || len(buf) == 0 {
			http.NotFound(w, r)
			log.Println("read json body failed", err)
			return
		}
		ret := make(map[string]string)
		err = json.Unmarshal(buf, &ret)
		if err != nil {
			http.NotFound(w, r)
			log.Println("Unmarshal json failed", err)
			return
		}
		for i, v := range ret {
			r.Form.Add(i, v)
			r.PostForm.Add(i, v)
		}
	} else {
		r.ParseForm()
	}
	w.Header().Add("connection", "close")
	if r.URL.EscapedPath() == "/keyen.do" {
		resp := make(map[string]string)
		key, _ := rsa.GenerateKey(rand.Reader, 2048)
		pkKey := key.Public().(*rsa.PublicKey)
		buf := x509.MarshalPKCS1PrivateKey(key)
		resp["privatekey"] = base64.StdEncoding.EncodeToString(buf)
		buf = x509.MarshalPKCS1PublicKey(pkKey)
		resp["pubkey"] = base64.StdEncoding.EncodeToString(buf)
		buf, _ = json.Marshal(resp)
		w.Write(buf)
		return
	} else if r.URL.EscapedPath() == "/call.do" {
		fc := r.FormValue("func")
		args := r.FormValue("args")
		sig := r.FormValue("signature")
		if len(fc) == 0 || len(args) == 0 || len(sig) == 0 {
			w.Write([]byte("invalid parameter"))
			return
		}
		obj := make(map[string]interface{})
		err := json.Unmarshal([]byte(args), &obj)
		if err != nil {
			w.Write([]byte("args is not json "))
			return
		}
		args_array := []string{fc, args, sig}
		w.Write([]byte(s.runner.SendRequest(args_array)))
		return
	} else {
		http.NotFound(w, r)
		return
	}
}

func (s *CoinServer) ServerHttp() {
	server := http.Server{}
	server.Addr = "0.0.0.0:8789"
	server.Handler = s
	server.ReadTimeout = 20 * time.Second
	server.WriteTimeout = 20 * time.Second
	server.ListenAndServe()
}

func main() {
	run, err := proxy.NewAppRunner()
	if err != nil {
		fmt.Println(err)
		return
	}
	s := &CoinServer{
		runner: run,
	}
	s.ServerHttp()
}

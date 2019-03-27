package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/xcoin/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
			http.Error(w, err.Error(), 500)
			fmt.Println("read json body failed", err)
			return
		}
		if buf == nil || len(buf) == 0 {
			http.Error(w, err.Error(), 500)
			fmt.Println("read json body failed", err)
			return
		}
		ret := make(map[string]string)
		err = json.Unmarshal(buf, &ret)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println("Unmarshal json failed", err)
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
	handler := proxy.NewAPIHandler(s.runner)
	handler.DispatchRequest(w, r)
}

func (s *CoinServer) ServerHttp() {
	server := http.Server{}
	server.Addr = "0.0.0.0:8789"
	server.Handler = s
	server.ReadTimeout = 2 * time.Minute
	server.WriteTimeout = 2 * time.Minute
	fmt.Println("working on :" + server.Addr)
	server.ListenAndServe()
}

func main() {
	confile := "./runner.conf"
	if len(os.Args) >= 2 {
		confile = os.Args[1]
	}
	run, err := proxy.NewAppRunner(confile)
	if err != nil {
		fmt.Println(err)
		return
	}
	s := &CoinServer{
		runner: run,
	}
	s.ServerHttp()
}

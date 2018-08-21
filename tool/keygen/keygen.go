package main

import (
	_ "crypto"
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	_ "encoding/json"
	"fmt"
	"os"
)

func useage() {
	fmt.Println("-g private key path public key path < format:ASN.1 DER>")
	//fmt.Println("-s srcfile private key <sign file>")
}

func genericKeyPair(private, public string) {
	fp, err := os.Create(private)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fp.Close()
	fpp, err := os.Create(public)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fpp.Close()
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	pkKey := key.Public().(*rsa.PublicKey)
	buf := x509.MarshalPKCS1PrivateKey(key)
	fp.WriteString(base64.StdEncoding.EncodeToString(buf))
	buf = x509.MarshalPKCS1PublicKey(pkKey)
	fpp.WriteString(base64.StdEncoding.EncodeToString(buf))
}

func main() {
	if len(os.Args) != 4 {
		useage()
		return
	}
	genericKeyPair(os.Args[2], os.Args[3])
}

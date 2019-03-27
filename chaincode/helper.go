package main

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"
)

func stringToPublicKey(text string) (interface{}, error) {
	buf, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKIXPublicKey(buf)
}

func publicKeyToString(pubk interface{}) (string, error) {
	buf, err := x509.MarshalPKIXPublicKey(pubk)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

func getSign(signature string) (rint, sint big.Int, err error) {
	byterun, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		err = errors.New("decrypt error," + err.Error())
		return
	}
	r, err := gzip.NewReader(bytes.NewBuffer(byterun))
	if err != nil {
		err = errors.New("decode error," + err.Error())
		return
	}
	defer r.Close()
	buf := make([]byte, 1024)
	count, err := r.Read(buf)
	if err != nil {
		err = errors.New("decode read error," + err.Error())
		return
	}
	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return
}

func checkSignature(arg string, sig string, pubkeyText string) error {
	hashed := sha256.Sum256([]byte(arg))
	pub, err := stringToPublicKey(pubkeyText)
	if err != nil {
		return err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		buf, err := base64.StdEncoding.DecodeString(sig)
		if err != nil {
			return err
		}
		return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], buf)
	case *ecdsa.PublicKey:
		r, s, err := getSign(sig)
		if err != nil {
			return err
		}
		if ecdsa.Verify(pub, hashed[:], &r, &s) {
			return nil
		}
		return errors.New("verify signature by ecdsa failed")
	default:
		return errors.New("unsupport signature  algorithm")
	}
}

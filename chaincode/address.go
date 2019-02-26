package main

import (
	"crypto/sha256"
	"math/big"
	"strings"
)

const (
	basePREFIX = "e1"
	alphabet   = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

var bigRadix = big.NewInt(58)
var bigZero = big.NewInt(0)

func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

// Base58Decode decode a modified base58 string to a byte slice, using alphabet
func Base58Decode(b string) []byte {
	return DecodeAlphabet(b, alphabet)
}

// Base58Encode a byte slice to a modified base58 string, using BTCAlphabet
func Base58Encode(b []byte) string {
	return EncodeAlphabet(b, alphabet)
}

func prefixAppend(b string) string {
	var tmp string
	tmp = basePREFIX + b
	return tmp
}

// EncodeWalletAddress encode wallet address from public key
func EncodeWalletAddress(publicKey []byte) string {
	hashed := sha256.Sum256(publicKey)
	tmp := EncodeAlphabet(hashed[:], alphabet)
	tmp = prefixAppend(tmp)
	return tmp
}

//IsWalletAddress check address
func IsWalletAddress(address string) bool {
	return strings.HasPrefix(address, basePREFIX) && len(address) < 100
}

// DecodeAlphabet decodes a modified base58 string to a byte slice, using alphabet.
func DecodeAlphabet(b, alphabet string) []byte {
	answer := big.NewInt(0)
	j := big.NewInt(1)
	for i := len(b) - 1; i >= 0; i-- {
		tmp := strings.IndexAny(alphabet, string(b[i]))
		if tmp == -1 {
			return []byte("")
		}
		idx := big.NewInt(int64(tmp))
		tmp1 := big.NewInt(0)
		tmp1.Mul(j, idx)

		answer.Add(answer, tmp1)
		j.Mul(j, bigRadix)
	}
	tmpval := answer.Bytes()
	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != alphabet[0] {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen, flen)
	copy(val[numZeros:], tmpval)
	return val
}

//EncodeAlphabet encodes a byte slice to a modified base58 string, using alphabet
func EncodeAlphabet(b []byte, alphabet string) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, alphabet[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabet[0])
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}

package utils

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/big"
	"math/rand"
	"time"
)

var letters = []byte("abcdefghjkmnpqrstuvwxyz123456789")
var longLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ=_")
var numLetters = []byte("123456789")

func init() {
	rand.Seed(time.Now().Unix())
}

func RandLowStr(n int) string {
	var _byte = RandLow(n)
	return string(_byte)
}

// RandNum 随机数字字符串，包含 1~9
func RandNum(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return ""
	}
	for i, x := range b {
		arc = x & 8
		b[i] = numLetters[arc]
	}
	return string(b)
}

// RandLow 随机字符串，包含 1~9 和 a~z - [i,l,o]
func RandLow(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	b := make([]byte, n)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return []byte{}
	}
	for i, x := range b {
		arc = x & 31
		b[i] = letters[arc]
	}
	return b
}

// RandUp 随机字符串，包含 英文字母和数字附加=_两个符号
func RandUp(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	b := make([]byte, n)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return []byte{}
	}
	for i, x := range b {
		arc = x & 63
		b[i] = longLetters[arc]
	}
	return b
}

// RandHex 生成16进制格式的随机字符串
func RandHex(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	var need int
	if n&1 == 0 { // even
		need = n
	} else { // odd
		need = n + 1
	}
	size := need / 2
	dst := make([]byte, need)
	src := dst[size:]
	if _, err := rand.Read(src[:]); err != nil {
		return []byte{}
	}
	hex.Encode(dst, src)
	return dst[:n]
}

func RandId(ids []int) (id int) {
	if len(ids) == 0 {
		return 0
	}
	bInt, _ := crand.Int(crand.Reader, big.NewInt(int64(len(ids))))
	id = int(bInt.Int64())
	return
}

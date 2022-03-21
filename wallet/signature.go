package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"runtime/debug"
)

type Signature struct {
	HexPrivateKey string
}

func NewSignature(hexPrivateKey string) Signature {
	//var hexPrivateKey = viper.GetString("sign_wallet.private_key")
	if len(hexPrivateKey) != 66 {
		panic(`privateKey: fail`)
	}
	return Signature{
		HexPrivateKey: hexPrivateKey,
	}
}

// sign 签名
func (s Signature) sign(data []byte) (signature string, err error) {
	var (
		privateKey    *ecdsa.PrivateKey
		signatureByte []byte
	)
	privateKey, err = crypto.HexToECDSA(s.HexPrivateKey[2:])
	signatureByte, err = crypto.Sign(data, privateKey)
	if err != nil {
		return
	}
	signatureByte[64] += 27
	signature = hexutil.Encode(signatureByte)
	return
}

// sha3Hash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calculated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func (s Signature) sha3Hash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func (s Signature) VerifySign(from, sigHex string, msg []byte) bool {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic recover! p: %v", p)
			debug.PrintStack()
		}
	}()
	fromAddr := common.HexToAddress(from)
	sig := hexutil.MustDecode(sigHex)
	if sig[64] != 27 && sig[64] != 28 {
	} else {
		sig[64] -= 27
	}
	pubKey, err := crypto.SigToPub(s.sha3Hash(msg), sig)
	if err != nil {
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return fromAddr == recoveredAddr
}

//SigToPub 验证签名 返回解析后的地址
func (s Signature) SigToPub(sigHex string, msg []byte) string {
	sig := hexutil.MustDecode(sigHex)
	if sig[64] != 27 && sig[64] != 28 {
	} else {
		sig[64] -= 27
	}
	pubKey, err := crypto.SigToPub(s.sha3Hash(msg), sig)
	if err != nil {
		log.Println("SigToPubErr:", err)
		return ""
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr.String()
}

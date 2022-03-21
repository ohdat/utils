package wallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"runtime/debug"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// TODO 本参考文件 不使用 随时会删除
const (
	//一番赏充值
	ichibanRechargeKey = "ichiban_recharge"
	//一番赏中心化提现
	ichibanWithdrawKey = "ichiban_withdraw"
	//币支付
	//一番赏充值
	tokenCollectKey = "token_collect"
	//1155质押
	stake1155Key = "stake_1155"
	//1155取回
	redeem1155Key = "redeem_1155"
)

type WalletService struct {
	HexPrivateKey string
}

func NewWalletService() WalletService {
	return WalletService{}
}

// sha3Hash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calculated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func (s WalletService) sha3Hash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

//SignPurchase  购买签名
func (s WalletService) SignPurchase(gameId int, address string, nonce string, amount int, payType int) (hashStr, signature string, err error) {
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(big.NewInt(int64(gameId)).Bytes(), 32),
		common.HexToAddress(address).Bytes(),
		common.LeftPadBytes(big.NewInt(int64(amount)).Bytes(), 32),
		[]byte(nonce),
		common.LeftPadBytes(big.NewInt(int64(payType)).Bytes(), 32),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//SignMortgage 兑换奖券/Eth 签名
func (s WalletService) SignMortgage(address string, nonce string, amount int, mortgageType string, tokenIds ...int) (hashStr, signature string, err error) {
	var tokenIdsBytes = make([][]byte, 0)
	for i := 0; i < len(tokenIds); i++ {
		tokenIdsBytes = append(tokenIdsBytes, common.LeftPadBytes(big.NewInt(int64(tokenIds[i])).Bytes(), 32))
	}
	tokenIdsBytes = append(tokenIdsBytes,
		common.HexToAddress(address).Bytes(),
		common.LeftPadBytes(big.NewInt(int64(amount)).Bytes(), 32),
		[]byte(nonce),
		[]byte(mortgageType),
	)

	hash := crypto.Keccak256Hash(
		tokenIdsBytes...,
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//SignWithdrawNFT 提现NFT 签名
func (s WalletService) SignWithdrawNFT(address string, nonce string, tokenId int, tokenAddress string, fromAddress string, fromTokenId int) (hashStr, signature string, err error) {
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(big.NewInt(int64(tokenId)).Bytes(), 32),
		common.HexToAddress(tokenAddress).Bytes(),
		common.LeftPadBytes(big.NewInt(int64(fromTokenId)).Bytes(), 32),
		common.HexToAddress(fromAddress).Bytes(),
		common.HexToAddress(address).Bytes(),
		[]byte(nonce),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//Erc20Stake Erc20 Stake UnStake 的签名
func (s WalletService) Erc20Stake(account big.Int, address string, nonce string, contract string, blockNum int, key string) (hashStr, signature string, err error) {
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(account.Bytes(), 32),
		common.HexToAddress(address).Bytes(),
		[]byte(nonce),
		common.HexToAddress(contract).Bytes(),
		common.LeftPadBytes(big.NewInt(int64(blockNum)).Bytes(), 32),
		[]byte(key),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//SignWithdrawMintNFT 提现NFT 签名
func (s WalletService) SignWithdrawMintNFT(address string, nonce string, tokenId int, tokenAddress string, fromTokenId int, creator string) (hashStr, signature string, err error) {
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(big.NewInt(int64(tokenId)).Bytes(), 32),
		common.HexToAddress(tokenAddress).Bytes(),
		common.LeftPadBytes(big.NewInt(int64(fromTokenId)).Bytes(), 32),
		common.HexToAddress(address).Bytes(),
		[]byte(nonce),
		common.HexToAddress(creator).Bytes(),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//MortgageBlackFriday 双十一抵押 签名
func (s WalletService) MortgageBlackFriday(address string, nonce string, tokenIds ...int) (hashStr, signature string, err error) {
	var tokenIdsBytes = make([][]byte, 0)
	for i := 0; i < len(tokenIds); i++ {
		tokenIdsBytes = append(tokenIdsBytes, common.LeftPadBytes(big.NewInt(int64(tokenIds[i])).Bytes(), 32))
	}
	tokenIdsBytes = append(tokenIdsBytes,
		common.HexToAddress(address).Bytes(),
		[]byte(nonce),
	)

	hash := crypto.Keccak256Hash(
		tokenIdsBytes...,
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//RecapBlackFriday 双十一抵押 签名
func (s WalletService) RecapBlackFriday(address string, nonce string) (hashStr, signature string, err error) {
	hash := crypto.Keccak256Hash(
		common.HexToAddress(address).Bytes(),
		[]byte(nonce),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//IchibanRechargeSign 一番赏中心化充值 签名
func (s WalletService) IchibanRechargeSign(eth string, ticket int, address string, nonce string) (hashStr, signature string, err error) {
	var ethBig *big.Int
	ethBig, err = s.String2BigInt(eth)
	if err != nil {
		return
	}
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(ethBig.Bytes(), 32),
		common.LeftPadBytes(big.NewInt(int64(ticket)).Bytes(), 32),
		common.HexToAddress(address).Bytes(),
		[]byte(nonce),
		[]byte(ichibanRechargeKey),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//NFT721AMint NFT721A铸造
func (s WalletService) NFT721AMint(quantity int, buyerAddress string, nonce string, tag string) (hashStr, signature string, err error) {
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(big.NewInt(int64(quantity)).Bytes(), 32),
		common.HexToAddress(buyerAddress).Bytes(),
		[]byte(nonce),
		[]byte(tag),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//Stake1155 质押1155
func (s WalletService) Stake1155(tokenIds []int, tokenAddress string, operatorAddress string, blockHeight int, nonce string) (hashStr, signature string, err error) {
	var tokenIdsBytes = make([][]byte, 0)
	for i := 0; i < len(tokenIds); i++ {
		tokenIdsBytes = append(tokenIdsBytes, common.LeftPadBytes(big.NewInt(int64(tokenIds[i])).Bytes(), 32))
	}
	tokenIdsBytes = append(tokenIdsBytes,
		common.HexToAddress(tokenAddress).Bytes(),
		common.HexToAddress(operatorAddress).Bytes(),
		common.LeftPadBytes(big.NewInt(int64(blockHeight)).Bytes(), 32),
		[]byte(nonce),
		[]byte(stake1155Key),
	)

	hash := crypto.Keccak256Hash(
		tokenIdsBytes...,
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//Redeem1155 取回1155
func (s WalletService) Redeem1155(tokenIds []int, tokenAddress string, operatorAddress string, blockHeight int, nonce string) (hashStr, signature string, err error) {
	var tokenIdsBytes = make([][]byte, 0)
	for i := 0; i < len(tokenIds); i++ {
		tokenIdsBytes = append(tokenIdsBytes, common.LeftPadBytes(big.NewInt(int64(tokenIds[i])).Bytes(), 32))
	}
	tokenIdsBytes = append(tokenIdsBytes,
		common.HexToAddress(tokenAddress).Bytes(),
		common.HexToAddress(operatorAddress).Bytes(),
		common.LeftPadBytes(big.NewInt(int64(blockHeight)).Bytes(), 32),
		[]byte(nonce),
		[]byte(redeem1155Key),
	)

	hash := crypto.Keccak256Hash(
		tokenIdsBytes...,
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

func (s WalletService) String2BigInt(tokenId string) (*big.Int, error) {
	n := new(big.Int)
	n, ok := n.SetString(tokenId, 10)
	if !ok {
		return nil, errors.New("SetString: error")
	}
	return n, nil
}

//IchibanWithdrawSign 一番赏中心化提现 签名
func (s WalletService) IchibanWithdrawSign(amount int, tokenId string, contractAddress, fromAddress, toAddress, nonce string) (hashStr, signature string, err error) {
	var tokenIdBig *big.Int
	tokenIdBig, err = s.String2BigInt(tokenId)
	if err != nil {
		return
	}
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(big.NewInt(int64(amount)).Bytes(), 32),
		common.LeftPadBytes(tokenIdBig.Bytes(), 32),
		common.HexToAddress(contractAddress).Bytes(),
		common.HexToAddress(fromAddress).Bytes(),
		common.HexToAddress(toAddress).Bytes(),
		[]byte(nonce),
		[]byte(ichibanWithdrawKey),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

//IchibanWithdrawBatchSign 一番赏中心化批量提现签名
func (s WalletService) IchibanWithdrawBatchSign(amounts []int, tokenIds []string, contractAddress, fromAddress []string, toAddress string, nonce string) (hashStr, signature string, err error) {

	var hashBytes = make([][]byte, 0)

	for i := 0; i < len(amounts); i++ {
		hashBytes = append(hashBytes, common.LeftPadBytes(big.NewInt(int64(amounts[i])).Bytes(), 32))
	}
	for i := 0; i < len(tokenIds); i++ {
		tokenIdBig, _ := s.String2BigInt(tokenIds[i])
		hashBytes = append(hashBytes, common.LeftPadBytes(tokenIdBig.Bytes(), 32))
	}

	for i := 0; i < len(contractAddress); i++ {
		hashBytes = append(hashBytes, common.HexToAddress(contractAddress[i]).Bytes())
	}

	for i := 0; i < len(fromAddress); i++ {
		hashBytes = append(hashBytes, common.HexToAddress(fromAddress[i]).Bytes())
	}

	hashBytes = append(hashBytes,
		common.HexToAddress(toAddress).Bytes(),
		[]byte(nonce),
		[]byte(ichibanWithdrawKey),
	)

	hash := crypto.Keccak256Hash(
		hashBytes...,
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.Sign(hashByte)
	return
}

// Sign 签名
func (s WalletService) Sign(data []byte) (signature string, err error) {
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

//SigToPub 验证签名 返回解析后的地址
func (s WalletService) SigToPub(sigHex, msg string) string {
	sig := hexutil.MustDecode(sigHex)
	if sig[64] != 27 && sig[64] != 28 {
	} else {
		sig[64] -= 27
	}
	pubKey, err := crypto.SigToPub(s.sha3Hash([]byte(msg)), sig)
	if err != nil {
		log.Println("SigToPubErr:", err)
		return ""
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	fmt.Println("32342")
	return recoveredAddr.String()
}

//VerifySign 验证签名
func (s WalletService) VerifySign(from, sigHex, msg string) bool {
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
	pubKey, err := crypto.SigToPub(s.sha3Hash([]byte(msg)), sig)
	if err != nil {
		log.Println("VerifySignErr", msg)
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return fromAddr == recoveredAddr
}

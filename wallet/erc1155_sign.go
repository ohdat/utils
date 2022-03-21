package wallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type Erc1155SignKey string

const (
	//Stake1155 1155质押
	Stake1155 Erc1155SignKey = "stake_1155"
	//Redeem1155 1155取回
	Redeem1155 Erc1155SignKey = "redeem_1155"
)

//Stake1155 质押1155
func (s Signature) Stake1155(
	tokenIds []int,
	tokenAddress string,
	operatorAddress string,
	blockHeight int,
	nonce string,
) (hashStr, signature string, err error) {
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
	signature, err = s.sign(hashByte)
	return
}

//Redeem1155 取回1155
func (s Signature) Redeem1155(
	tokenIds []int,
	tokenAddress string,
	operatorAddress string,
	blockHeight int,
	nonce string,
) (hashStr, signature string, err error) {
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
	signature, err = s.sign(hashByte)
	return
}

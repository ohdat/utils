package wallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type Erc20SignKey string

const (
	StakeToken   Erc20SignKey = "stake_token"
	RedeemToken  Erc20SignKey = "redeem_token"
	TokenCollect Erc20SignKey = "token_collect"
	TokenPay     Erc20SignKey = "token_pay"
)

//Erc20Stake Erc20 Stake Redeem 的签名
func (s Signature) Erc20Stake(
	amount big.Int,
	address, contract string,
	blockNum int,
	nonce string,
	key Erc20SignKey,
) (hashStr, signature string, err error) {
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(amount.Bytes(), 32),
		common.HexToAddress(address).Bytes(),
		common.HexToAddress(contract).Bytes(),
		common.LeftPadBytes(big.NewInt(int64(blockNum)).Bytes(), 32),
		[]byte(nonce),
		[]byte(key),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.sign(hashByte)
	return
}

//TokenCollect 签名
/**
function hashTokenCollect(
	address _operatorAddress,
	uint256 _amount,
	address _tokenAddress,
	address _fromAddress,
	uint256 _blockHeight,
	string memory _nonce
) private pure returns (bytes32) {
	bytes32 hash = keccak256(
		abi.encodePacked(
			"\x19Ethereum Signed Message:\n32",
			keccak256(
				abi.encodePacked(
					_operatorAddress,
					_amount,
					_tokenAddress,
					_fromAddress,
					_blockHeight,
					_nonce,
					"token_collect"
				)
			)
		)
	);
	return hash;
}
*/
func (s Signature) TokenCollect(
	amount big.Int,
	address,
	tokenAddress,
	formAddress string,
	blockNum int,
	nonce string) (hashStr, signature string, err error) {
	hash := crypto.Keccak256Hash(
		common.HexToAddress(address).Bytes(),
		common.LeftPadBytes(amount.Bytes(), 32),
		common.HexToAddress(tokenAddress).Bytes(),
		common.HexToAddress(formAddress).Bytes(),
		common.LeftPadBytes(big.NewInt(int64(blockNum)).Bytes(), 32),
		[]byte(nonce),
		[]byte(TokenCollect),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.sign(hashByte)
	return
}

//TokenPay 20币支付 签名
func (s Signature) TokenPay(
	operatorAddress string,
	amount *big.Int,
	tokenAddress string,
	toAddress string,
	divideAmount *big.Int,
	blockHeight int,
	nonce string,
) (hashStr, signature string, err error) {

	hash := crypto.Keccak256Hash(
		common.HexToAddress(operatorAddress).Bytes(),
		common.LeftPadBytes(amount.Bytes(), 32),
		common.HexToAddress(tokenAddress).Bytes(),
		common.HexToAddress(toAddress).Bytes(),
		common.LeftPadBytes(divideAmount.Bytes(), 32),
		common.LeftPadBytes(big.NewInt(int64(blockHeight)).Bytes(), 32),
		[]byte(nonce),
		[]byte(TokenPay),
	)
	hashByte := s.sha3Hash(hash.Bytes())
	hashStr = hexutil.Encode(hashByte)
	signature, err = s.sign(hashByte)
	return
}

package wallet

import (
	"fmt"
	"github.com/mytokenio/ethrpc"
	"log"
	"time"
)

type Ether struct {
	EthRpc    ethrpc.EthRPC
	StepBlock int //签名过期高度

}

func NewEther(rpcUri string, stepBlock int) Ether {
	if rpcUri == "" {
		log.Fatalf("rpcUri is nil")
	}
	if stepBlock < 1 {
		stepBlock = 20
	}
	return Ether{
		EthRpc:    ethrpc.NewNodeAPI(rpcUri),
		StepBlock: stepBlock,
	}
}

//BlockNum 获取以太坊当前高度
func (s Ether) BlockNum() int {
	block, err := s.EthRpc.EthBlockNumber()

	if err != nil {
		fmt.Printf("getBlockNumError: %v \n", err)
		time.Sleep(time.Millisecond * 200) //200 毫秒
		return s.BlockNum()
	}
	return block
}

//ExpiredBlock 获取以太坊当前高度 + 过期高度
func (s Ether) ExpiredBlock() (block int) {
	block = s.BlockNum()
	block = block + s.StepBlock
	return
}

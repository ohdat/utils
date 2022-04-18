package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestStructAssign(t *testing.T) {
	//fmt.Println("price", price)
	type S1 struct {
		SSS string
		DDD string
	}

	type S2 struct {
		SSS int
		DDD string
	}

	var s1 = &S1{"123", "123"}
	var s2 = new(S2)

	StructAssign(s2, s1)
	fmt.Println("s2:", s2)
}

func NewRedisClient() *redis.Client {
	opts := redis.Options{
		Addr:         "54.180.12.56:2379",
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
	}

	client := redis.NewClient(&opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic("ping redis error")
	}

	return client
}

func TestBloomFilter(t *testing.T) {
	client := NewRedisClient()
	if notExist, err := bloomFilterNotExist(context.Background(), client, "bf_key", "blocknumber"); err != nil {
		fmt.Printf("bf_key occur error:%s", err)
	} else {
		fmt.Printf("bf_key:%t", notExist)
	}
}

func TestBloomFilterByBitSet(t *testing.T) {
	client := NewRedisClient()

	//use the estimate m and k.
	m, k := EstimateParameters(100000, 0.001)
	// implement BitSetProvider
	bitSet := NewRedisBitSet("test_key", m, client)
	//new a Bloom Filter
	bf := NewBloomFilter(m, k, bitSet)

	for i := 0; i < 5000; i++ {
		fmt.Println(i)
		//check exist
		data := []byte(fmt.Sprintf("%s%d", "hello", i))
		if notExists, err := bf.Test(data); err != nil {
			fmt.Sprintf("%s,%v\n", "判断出错了", err)

		} else {
			if notExists {
				fmt.Println("一定不存在data", notExists, err)
				errAdd := bf.Add(data)
				fmt.Sprintf("%v\n", errAdd)
			} else {
				fmt.Println("可能存在data", notExists, err)
			}
		}
	}

}

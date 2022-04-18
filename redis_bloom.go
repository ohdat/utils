package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"hash/fnv"
	"math"
)

// 使用布隆过滤器，如果不存在返回true。如果返回false，不一定存在。只能用于判断不存在，不能用于判断存在
// RedisBloom需要单独安装才能使用这些命令
func bloomFilterNotExist(ctx context.Context, c *redis.Client, key string, value string) (bool, error) {
	exists, err := c.Do(ctx, "BF.EXISTS", key, value).Bool()
	if err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}

	return false, nil
}

type BitSetProvider interface {
	Set([]int64) error
	Test([]int64) (bool, error)
}

type BloomFilter struct {
	m      int64 // bloomfilter's  size (bit)
	k      int64 // hash function count
	bitset BitSetProvider
}

func NewBloomFilter(m int64, k int64, bitset BitSetProvider) *BloomFilter {
	return &BloomFilter{
		m:      m,
		k:      k,
		bitset: bitset,
	}
}

// n 是总量
// p 错误率
// m size(bit) k hash函数的个数
// https://krisives.github.io/bloom-calculator/
func EstimateParameters(n uint, p float64) (int64, int64) {
	m := math.Ceil(float64(n) * math.Log(p) / math.Log(1.0/math.Pow(2.0, math.Ln2)))
	k := math.Ln2*m/float64(n) + 0.5

	return int64(m), int64(k)
}

func (f *BloomFilter) Add(data []byte) error {
	location := f.getLocations(data)

	err := f.bitset.Set(location)

	fmt.Println("add", location, err)

	if err != nil {
		return err
	}
	return nil
}

func (f *BloomFilter) Test(data []byte) (bool, error) {
	location := f.getLocations(data)
	notExist, err := f.bitset.Test(location)

	if err != nil {
		return false, err
	}

	if !notExist {
		fmt.Println("一定不存在data", notExist, err)
		// 这里返回一定不存在
		return true, nil
	}

	// 如果存在，原值data不一定存在
	return false, nil
}

func (f *BloomFilter) getLocations(data []byte) []int64 {
	locations := make([]int64, f.k)
	hasher := fnv.New64()
	hasher.Write(data)
	a := make([]byte, 1)
	for i := int64(0); i < f.k; i++ {
		a[0] = byte(i)
		hasher.Write(a)
		hashValue := hasher.Sum64()
		locations[i] = int64(hashValue % uint64(f.m))
	}
	return locations
}

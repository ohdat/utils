package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

// redis_bitset
const redisMaxLength int64 = 8 * 512 * 1024 * 1024 //512M

type RedisBitSet struct {
	keyPrefix string
	m         int64
	client    *redis.Client
}

func NewRedisBitSet(keyPrefix string, m int64, client *redis.Client) *RedisBitSet {
	return &RedisBitSet{keyPrefix, m, client}
}

/*
//sync set
func (r *RedisBitSet) Set(offsets []int64) error {
	for _, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		_, err := r.client.SetBit(key, thisOffset, 1).Result()
		if err != nil {
			return err
		}
	}
	return nil
}
*/

//async set
func (r *RedisBitSet) Set(offsets []int64) error {
	var wg sync.WaitGroup
	doneChan := make(chan struct{}, 1)
	errChan := make(chan error, 1)
	for _, off := range offsets {
		wg.Add(1)
		go func(offset int64) {
			defer wg.Done()
			key, thisOffset := r.getKeyOffset(offset)
			_, err := r.client.SetBit(context.Background(), key, thisOffset, 1).Result()
			if err != nil {
				errChan <- err
				return
			}
		}(off)
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case <-doneChan:
		return nil
	case err := <-errChan:
		return err
	}
}

func (r *RedisBitSet) Test(offsets []int64) (bool, error) {
	for _, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		bitValue, err := r.client.GetBit(context.Background(), key, thisOffset).Result()
		if err != nil {
			return false, err
		}
		if bitValue == 0 {
			return false, nil
		}
	}

	return true, nil
}

// Set Expire time if needed
func (r *RedisBitSet) Expire(seconds uint) error {
	max := int(r.m / redisMaxLength)

	for n := 0; n <= max; n++ {
		key := fmt.Sprintf("%s:%d", r.keyPrefix, n)
		_, err := r.client.Expire(context.Background(), key, time.Duration(seconds)*time.Second).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete all the keys
func (r *RedisBitSet) Delete() error {
	max := int(r.m / redisMaxLength)
	keys := make([]string, 0)
	for n := 0; n <= max; n++ {
		key := fmt.Sprintf("%s:%d", r.keyPrefix, n)
		keys = append(keys, key)
		n = n + 1
	}
	_, err := r.client.Del(context.Background(), keys...).Result()
	return err
}

func (r *RedisBitSet) getKeyOffset(offset int64) (string, int64) {
	index := int64(offset / redisMaxLength)
	thisOffset := offset - index*redisMaxLength

	key := fmt.Sprintf("%s:%d", r.keyPrefix, index)

	return key, thisOffset
}

var _ BitSetProvider = (*RedisBitSet)(nil)

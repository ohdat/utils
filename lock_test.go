package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"testing"
	"time"
)

func TestGetLock(t *testing.T) {
	viper.SetConfigFile("config.yaml")
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	opts := redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.durable_db"),
		PoolSize:     viper.GetInt("redis.pool_size"),
		MinIdleConns: viper.GetInt("redis.min_idle_conns"),
	}

	client := redis.NewClient(&opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("PersistenceRedis 客户端初始化错误 %v", err.Error()))
	}
	fmt.Println("PersistenceRedis 客户端初始化完成")
	go func() {
		lock := GetLock(client, "sfsdf")
		err := lock.Lock()
		if err != nil {
			log.Fatalln("GetLock err", err)
		}
		defer lock.Unlock()
		fmt.Println("11111:",time.Now())
		time.Sleep(time.Second * 2)
	}()
	go func() {
		lock := GetLock(client, "sfsdf")
		err := lock.Lock()
		if err != nil {
			log.Fatalln("GetLock err", err)
		}
		defer lock.Unlock()
		fmt.Println("2222:",time.Now())
		time.Sleep(time.Second * 2)
	}()
	time.Sleep(time.Second * 4)
}

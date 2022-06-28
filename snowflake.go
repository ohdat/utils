package utils

import (
	"github.com/sony/sonyflake"
	"sync"
	"time"
)

var sf *sonyflake.Sonyflake
var sfOnce sync.Once

func GetSnowflake() *sonyflake.Sonyflake {
	sfOnce.Do(func() {
		sf = sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: time.Time{},
		})
	})
	return sf
}

func NextSnowflakeID() (id uint64) {
	id, err := GetSnowflake().NextID()
	if err != nil {
		panic(err)
	}
	return id
}

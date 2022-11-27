package store

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var (
	redisClient     *redis.Client
	redisClientOnce sync.Once
)

func Redis() *redis.Client {
	redisClientOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d",
				Conf.Host,
				Conf.Port,
			),
			Password:     Conf.Password,
			DB:           Conf.Database,
			PoolSize:     15,
			MinIdleConns: 10,
			DialTimeout:  5 * time.Second, // 超时时间
			ReadTimeout:  3 * time.Second, // 读取超时时间
		})
		pong, err := redisClient.Ping(redisClient.Context()).Result()
		if err != nil {
			panic(fmt.Sprintf("redis初始化连接失败：%s", err.Error()))
		}
		fmt.Println("redis链接"+pong, err)
	})
	return redisClient
}

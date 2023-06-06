/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  redis.go  redis.go 2022-11-30
 */

package store

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type RedisStore struct {
	client *redis.Client
	Once   sync.Once
}

type RedisConfig struct {
	Host     string
	Port     int
	Prefix   string
	Password string
	Database int
}

var (
	redisStore RedisStore
)

func Redis(c RedisConfig) *redis.Client {
	redisStore.Once.Do(func() {
		redisStore.client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d",
				c.Host,
				c.Port,
			),
			Password:     c.Password,
			DB:           c.Database,
			PoolSize:     15,
			MinIdleConns: 10,
			DialTimeout:  5 * time.Second, // 超时时间
			ReadTimeout:  3 * time.Second, // 读取超时时间
			// 开启 notify-keyspace-events KEA
		})
		ctx := context.Background()
		pong, err := redisStore.client.Ping(ctx).Result()
		if err != nil {
			panic(fmt.Sprintf("redis初始化连接失败：%s", err.Error()))
			return
		}
		fmt.Println("redis链接"+pong, err)
	})
	return redisStore.client
}

// GetClient 获取client
func (r *RedisStore) GetClient() *redis.Client {
	return r.client
}

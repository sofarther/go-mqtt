package initialize

import (
	"github.com/go-redis/redis/v7"
)

var client *redis.Client

const address  = "127.0.0.1:6379"
const password = "publink"
const db = 0

func init()  {
	Redis()
}
func Redis() *redis.Client {
    // 需加锁 判断，保证 client 只初始化 一次
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr: address,
			Password: password,
			DB: db,

		})

	}

	return client
}

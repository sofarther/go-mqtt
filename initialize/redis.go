package initialize

import (
	"github.com/go-redis/redis/v7"
)

var client *redis.Client

const address  = "127.0.0.1:6379"
const password = "publink"
const db = 0


func Redis() *redis.Client {

	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr: address,
			Password: password,
			DB: db,

		})

	}

	return client
}

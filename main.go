package main

import (
	"log"
	"proxy_project/cache/redis"
	"proxy_project/proxyAPI"
)

func main() {
	redisClient := redis.RedisClient{}
	err := redisClient.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}

	proxyAPI.ProxyAPI(&redisClient)
}

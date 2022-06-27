package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"proxy_project/api"
	"proxy_project/cache/redis"
	"proxy_project/proxy"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisClient := redis.RedisClient{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	requestClient := api.RequestClient{
		Host: os.Getenv("API_HOST"),
		Port: os.Getenv("API_PORT"),
	}

	err = redisClient.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}

	proxyAPI := proxy.ProxyAPI{
		CacheClient:   &redisClient,
		RequestClient: requestClient,
	}

	if err := proxyAPI.Run(); err != nil {
		log.Fatal(err)
	}
}

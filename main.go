package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"proxy_project/cache/redis"
	"proxy_project/proxyAPI"
	"proxy_project/proxyAPI/requestAPI"
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

	requestCLient := requestAPI.RequestToAPIClient{
		Host: os.Getenv("API_HOST"),
		Port: os.Getenv("API_PORT"),
	}

	err = redisClient.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}

	proxyAPI.ProxyAPI(&redisClient, &requestCLient)
}

func aaa() { fmt.Println("seila") }

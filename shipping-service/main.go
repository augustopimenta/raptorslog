package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

type Order struct {
	ID       string `json:"id"`
	Version  int    `json:"version"`
	Location string `json:"location"`
}

var (
	client *redis.Client
)

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     getEnv("QUEUE_HOST", "localhost:6379"),
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()

	fmt.Println(pong, err)

	for {
		result, _ := client.BLPop(0, "queue:orders").Result()

		fmt.Println("Processing:", result[1])
	}
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	return value
}

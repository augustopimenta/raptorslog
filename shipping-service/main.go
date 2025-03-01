package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

type Order struct {
	ID       string `json:"id"`
	Version  int    `json:"version"`
	Location string `json:"location"`
}

const (
	QueueHostEnv   = "QUEUE_HOST"
	TruckAMHostEnv = "TRUCK_AM_HOST"
	TruckMGHostEnv = "TRUCK_MG_HOST"
	TruckRSHostEnv = "TRUCK_RS_HOST"
)

var (
	client         *redis.Client
	deliveryRoutes map[string]string
)

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     getEnv(QueueHostEnv, "localhost:6379"),
		Password: "",
		DB:       0,
	})

	fillRoutes()

	fmt.Println("shipping-service:v2 started!")

	ctrl := make(chan int, 100)

	for {
		ctrl <- 1

		var order Order
		result, _ := client.BLPop(0, "queue:orders").Result()

		json.Unmarshal([]byte(result[1]), &order)

		if host, ok := deliveryRoutes[order.Location]; ok && host != "" {
			go func() {
				processOrder(host, result[1])
				<-ctrl
			}()
		} else {
			fmt.Printf("[v2] Processing %s, with error\n", result[1])
		}
	}
}

func processOrder(host, data string) {
	fmt.Printf("[v2] Processing: %s, sending to host %s\n", data, host)

	http.Post(fmt.Sprintf("http://%s/deliver", host), "application/json", bytes.NewBuffer([]byte(data)))
}

func fillRoutes() {
	deliveryRoutes = make(map[string]string)

	deliveryRoutes["AM"] = getEnv(TruckAMHostEnv, "")
	deliveryRoutes["MG"] = getEnv(TruckMGHostEnv, "")
	deliveryRoutes["RS"] = getEnv(TruckRSHostEnv, "")
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	return value
}

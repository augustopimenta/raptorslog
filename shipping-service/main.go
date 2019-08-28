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
	QueueHostEnv = "QUEUE_HOST"
	TruckAMHost  = "TRUCK_AM_HOST"
	TruckMGHost  = "TRUCK_MG_HOST"
	TruckRSHost  = "TRUCK_RS_HOST"
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

	fmt.Printf("shipping-service started! %v\n", deliveryRoutes)

	for {
		var order Order
		result, _ := client.BLPop(0, "queue:orders").Result()

		json.Unmarshal([]byte(result[1]), &order)

		if host, ok := deliveryRoutes[order.Location]; ok && host != "" {
			fmt.Printf("Processing: %s, sending to host %s\n", result[1], host)

			http.Post(fmt.Sprintf("http://%s/deliver", host), "application/json", bytes.NewBuffer([]byte(result[1])))
		} else {
			fmt.Printf("Processing %s, with error \n", result[1])
		}
	}
}

func fillRoutes() {
	deliveryRoutes = make(map[string]string)

	deliveryRoutes["AM"] = getEnv(TruckAMHost, "")
	deliveryRoutes["MG"] = getEnv(TruckMGHost, "")
	deliveryRoutes["RS"] = getEnv(TruckRSHost, "")
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	return value
}

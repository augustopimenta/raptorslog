package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Order struct {
	ID       string `json:"id"`
	Version  int    `json:"version"`
	Location string `json:"location"`
}

const (
	QueueHostEnv = "QUEUE_HOST"
)

var (
	locations = []string{"AM", "MG", "RS"}
	client    *redis.Client
)

func orderHandler(w http.ResponseWriter, r *http.Request) {
	order := Order{uuid.New().String(), 2, getRandomLocation()}

	message := fmt.Sprintf("[v%d] Order %s to %s\n", order.Version, order.ID, order.Location)

	json, _ := json.Marshal(order)

	client.RPush("queue:orders", json)

	log.Println(message)

	w.Write([]byte(message))
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     getEnv(QueueHostEnv, "localhost:6379"),
		Password: "",
		DB:       0,
	})

	fmt.Println("order-service:v2 started!")

	r := mux.NewRouter()
	r.HandleFunc("/order", orderHandler)

	log.Fatal(http.ListenAndServe(":80", r))
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	return value
}

func getRandomLocation() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(locations))

	return locations[i]
}

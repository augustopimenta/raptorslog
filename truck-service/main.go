package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Request struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}

var client *mongo.Client

func main() {
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://admin:admin@localhost/?w=majority"))

	r := mux.NewRouter()
	r.HandleFunc("/deliver", deliverHandler)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func deliverHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	json.NewDecoder(r.Body).Decode(&request)

	t := getDeliveryTime()

	<-time.After(time.Duration(t) * time.Second)

	collection := client.Database("raptorslog").Collection("deliveries")

	collection.InsertOne(context.TODO(), bson.M{"id": request.ID, "location": request.Location, "time": t})

	fmt.Printf("Order %s to %s delivered after %d seconds", request.ID, request.Location, t)
}

func getDeliveryTime() int {
	time, _ := strconv.Atoi(getEnv("DELIVERY_TIME", "0"))

	return time
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	return value
}

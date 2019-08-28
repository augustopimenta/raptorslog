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

const (
	DatabaseConnectionEnv = "DATABASE_CONNECTION"
	DatabaseNameEnv       = "DATABASE_NAME"
	DeliveryTimeEnv       = "DELIVERY_TIME"
)

var client *mongo.Client

func main() {
	dbConnection := getEnv(DatabaseConnectionEnv, "mongodb://admin:admin@localhost/?w=majority")

	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(dbConnection))

	fmt.Printf("shipping-service started! DELIVERY_TIME=%d\n", getDeliveryTime())

	r := mux.NewRouter()
	r.HandleFunc("/deliver", deliverHandler)

	log.Fatal(http.ListenAndServe(":80", r))
}

func deliverHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	json.NewDecoder(r.Body).Decode(&request)

	t := getDeliveryTime()

	<-time.After(time.Duration(t) * time.Second)

	collection := client.Database(getDatabaseName()).Collection("deliveries")

	collection.InsertOne(context.TODO(), bson.M{"id": request.ID, "location": request.Location, "time": t})

	fmt.Printf("Order %s to %s delivered after %d seconds", request.ID, request.Location, t)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getDatabaseName() string {
	return getEnv(DatabaseNameEnv, "raptorslog")
}

func getDeliveryTime() int {
	time, _ := strconv.Atoi(getEnv(DeliveryTimeEnv, "0"))

	return time
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	return value
}

package data

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var err error

type Order struct {
	OrderId     string            `bson:"order_id"`
	Items       map[string]string `bson:"items"`
	TotalAmount int               `bson:"total_amount"`
}

func init() {
	password := os.Getenv("PASS")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://user:" + password + "@cluster0.8ltsy.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Println("Error connecting to client:", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Printf("Error connecting to mongo server: %v", err)
	}

	collection = client.Database("ecommerce").Collection("order")
}

func (o *Order) Store() {
	_, err = collection.InsertOne(context.Background(), o)
	if err != nil {
		fmt.Println("Error inserting order:", err)
		return
	}
}

func (o *Order) Get() {
	filter := bson.D{{"order_id", "1"}}
	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(o); err != nil {
		fmt.Println("Error retreving order:", err)
	}
}

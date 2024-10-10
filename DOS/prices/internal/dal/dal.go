package dal

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"blaucorp.com/prices/internal/misc"
)

// Connect to MongoDB
func ConnectMongoDB() (*mongo.Client, context.Context) {
	clientOptions := options.Client().ApplyURI("mongodb://user:123qwe@localhost:27017/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	return client, ctx
}

// Function to create a price list
func CreatePriceList(db *mongo.Database, listName, owner string) {
	priceListCollection := db.Collection("price_lists")
	priceList := bson.D{
		{"list", listName},
		{"owner", owner},
	}
	_, err := priceListCollection.InsertOne(context.TODO(), priceList)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created price list '%s' for owner '%s'\n", listName, owner)
}

// Function to assign targets to the price list
func AssignTargets(db *mongo.Database, listName string, targets []string) {
	targetCollection := db.Collection("targets")
	for _, target := range targets {
		targetData := bson.D{
			{"list", listName},
			{"target", target},
		}
		_, err := targetCollection.InsertOne(context.TODO(), targetData)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Assigned targets %v to list '%s'\n", targets, listName)
}

// Function to add a price to the list
func AddPrice(db *mongo.Database, listName, sku, unit, material, tservicio string, price float64) {
	priceCollection := db.Collection("prices")
	priceTuple := map[string]string{
		"list":      listName,
		"sku":       sku,
		"unit":      unit,
		"material":  material,
		"tservicio": tservicio,
	}
	priceHash := misc.GenerateHash(priceTuple)
	priceData := bson.D{
		{"tuple", priceTuple},
		{"hash", priceHash},
		{"price", price},
	}
	_, err := priceCollection.InsertOne(context.TODO(), priceData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added price %.2f for SKU %s to list '%s'\n", price, sku, listName)
}

// Function to retrieve a price by tuple
func RetrievePriceByTuple(db *mongo.Database, priceTuple map[string]string) {
	priceCollection := db.Collection("prices")
	priceHash := misc.GenerateHash(priceTuple)

	filter := bson.D{{"hash", priceHash}}
	var result bson.M
	err := priceCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println("No price found for the given tuple.")
		return
	}

	fmt.Printf("Price found: SKU: %s, Price: %.2f, Hash: %s\n", result["tuple"].(bson.M)["sku"], result["price"], result["hash"])
}

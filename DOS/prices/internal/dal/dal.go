package dal

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"blaucorp.com/prices/internal/misc"
)

// Function to create a price list
func CreatePriceList(db *mongo.Database, listName, owner string) error {
	priceListCollection := db.Collection("price_lists")
	priceList := bson.D{
		{"list", listName},
		{"owner", owner},
	}
	_, err := priceListCollection.InsertOne(context.TODO(), priceList)
	if err != nil {
		return err
	}

	return nil
}

// Function to assign targets to the price list
func AssignTargets(db *mongo.Database, listName string, targets []string) error {

	targetCollection := db.Collection("targets")
	for _, target := range targets {
		targetData := bson.D{
			{"list", listName},
			{"target", target},
		}
		_, err := targetCollection.InsertOne(context.TODO(), targetData)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Assigned targets %v to list '%s'\n", targets, listName)
	return nil
}

// Function to add a price to the list
func AddPrice(db *mongo.Database, listName, sku, unit, material, tservicio string, price float64) error {
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
		return err
	}

	return nil
}

// Function to retrieve a price by tuple
func RetrievePriceByTuple(db *mongo.Database, priceTuple map[string]string) (float64, error) {
	priceCollection := db.Collection("prices")
	priceHash := misc.GenerateHash(priceTuple)

	filter := bson.D{{"hash", priceHash}}
	var result bson.M
	err := priceCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return 0, err
	}

	price, ok := result["price"].(float64)
	if !ok {
		return 0, fmt.Errorf("price not found or invalid type")
	}

	return price, nil
}

// Deletes a list along with its associated targets and prices.
func DeleteList(db *mongo.Database, listName string) error {
	ctx := context.TODO()

	// Delete associated targets
	_, err := db.Collection("targets").DeleteMany(ctx, bson.M{"list": listName})
	if err != nil {
		return fmt.Errorf("failed to delete targets for list '%s': %v", listName, err)
	}

	// Delete associated prices
	_, err = db.Collection("prices").DeleteMany(ctx, bson.M{"tuple.list": listName})
	if err != nil {
		return fmt.Errorf("failed to delete prices for list '%s': %v", listName, err)
	}

	// Delete the list itself
	_, err = db.Collection("lists").DeleteOne(ctx, bson.M{"list": listName})
	if err != nil {
		return fmt.Errorf("failed to delete list '%s': %v", listName, err)
	}

	return nil
}

package dal

import (
	"context"
	"errors"
	"fmt"
	"log"

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

	return nil
}

// Function to add a price to the list
func AddPrice(db *mongo.Database, listName, sku, unit, material, tservicio string, price float64) error {
	// Verify price list existence
	if !ExistsPriceList(db, listName) {
		return fmt.Errorf("Price list %s does not exist", listName)
	}

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

func EditPrice(db *mongo.Database, listName, sku, unit, material, tservicio string, price float64) error {
	// Verify price list existence
	if !ExistsPriceList(db, listName) {
		return fmt.Errorf("Price list %s does not exist", listName)
	}

	priceCollection := db.Collection("prices")
	priceTuple := map[string]string{
		"list":      listName,
		"sku":       sku,
		"unit":      unit,
		"material":  material,
		"tservicio": tservicio,
	}
	priceHash := misc.GenerateHash(priceTuple)

	// Define the filter to match the price based on the hash
	filter := bson.D{{"hash", priceHash}}

	// Check if the document exists
	var existingDoc bson.M
	err := priceCollection.FindOne(context.TODO(), filter).Decode(&existingDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Document does not exist, handle the error (or insert if needed)
			return errors.New("price document not found")
		}
		return err
	}

	// Document exists, perform the update
	update := bson.D{
		{"$set", bson.D{
			{"tuple", priceTuple},
			{"hash", priceHash},
			{"price", price},
		}},
	}

	_, err = priceCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// ClonePriceList clones a price list, its targets, and prices with a new name.
func ClonePriceList(db *mongo.Database, originalListName, newListName string) error {
	ctx := context.TODO()

	// Collections
	priceListCollection := db.Collection("price_lists")
	targetCollection := db.Collection("targets")
	priceCollection := db.Collection("prices")

	// Clone the price list
	var originalList bson.M
	err := priceListCollection.FindOne(ctx, bson.M{"list": originalListName}).Decode(&originalList)
	if err != nil {
		return fmt.Errorf("failed to find the original list: %v", err)
	}

	// Create the new price list with the new name
	originalList["list"] = newListName
	delete(originalList, "_id")
	_, err = priceListCollection.InsertOne(ctx, originalList)
	if err != nil {
		return fmt.Errorf("failed to insert the cloned list: %v", err)
	}
	fmt.Printf("Cloned price list '%s' as '%s'\n", originalListName, newListName)

	// Clone the targets associated with the original list
	cursor, err := targetCollection.Find(ctx, bson.M{"list": originalListName})
	if err != nil {
		return fmt.Errorf("failed to find targets: %v", err)
	}
	defer cursor.Close(ctx)

	var clonedTargets []interface{}
	for cursor.Next(ctx) {
		var target bson.M
		err := cursor.Decode(&target)
		if err != nil {
			log.Printf("failed to decode target: %v", err)
			continue
		}
		target["list"] = newListName // Update the list name
		delete(target, "_id")
		clonedTargets = append(clonedTargets, target)
	}
	if len(clonedTargets) > 0 {
		_, err = targetCollection.InsertMany(ctx, clonedTargets)
		if err != nil {
			return fmt.Errorf("failed to insert cloned targets: %v", err)
		}
		fmt.Printf("Cloned targets for list '%s'\n", newListName)
	}

	// Clone the prices associated with the original list
	cursor, err = priceCollection.Find(ctx, bson.M{"tuple.list": originalListName})
	if err != nil {
		return fmt.Errorf("failed to find prices: %v", err)
	}
	defer cursor.Close(ctx)

	var clonedPrices []interface{}
	for cursor.Next(ctx) {
		var price bson.M
		err := cursor.Decode(&price)
		if err != nil {
			log.Printf("failed to decode price: %v", err)
			continue
		}
		price["tuple"].(bson.M)["list"] = newListName // Update the list name in the tuple
		tpl := map[string]string{
			"list":      price["tuple"].(bson.M)["list"].(string),
			"sku":       price["tuple"].(bson.M)["sku"].(string),
			"unit":      price["tuple"].(bson.M)["unit"].(string),
			"material":  price["tuple"].(bson.M)["material"].(string),
			"tservicio": price["tuple"].(bson.M)["tservicio"].(string),
		}
		price["hash"] = misc.GenerateHash(tpl)
		delete(price, "_id")
		clonedPrices = append(clonedPrices, price)
	}
	if len(clonedPrices) > 0 {
		_, err = priceCollection.InsertMany(ctx, clonedPrices)
		if err != nil {
			return fmt.Errorf("failed to insert cloned prices: %v", err)
		}
		fmt.Printf("Cloned prices for list '%s'\n", newListName)
	}

	return nil
}

func ExistsPriceList(db *mongo.Database, listName string) bool {

	priceListCollection := db.Collection("price_lists")
	var originalList bson.M

	err := priceListCollection.FindOne(context.TODO(), bson.M{"list": listName}).Decode(&originalList)
	if err != nil {
		return false
	}
	return true
}

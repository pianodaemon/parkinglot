package dal

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"blaucorp.com/prices/internal/misc"
)

// Function to retrieve a price by tuple
func RetrievePriceByTuple(db *mongo.Database, priceTuple map[string]string) (float64, error) {
	priceCollection := db.Collection("prices")
	priceHash := misc.GenerateHash(priceTuple)

	filter := bson.D{{"hash", priceHash}}
	var result bson.M
	err := priceCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("price not found")
	}

	price, ok := result["price"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid type")
	}

	return price, nil
}

// Retrieves all lists associated with a given owner and specified targets.
func GetListsByOwnerAndTargets(db *mongo.Database, owner string, targets []string) ([]string, error) {
	ctx := context.TODO()
	listCollection := db.Collection("price_lists")

	// Find lists associated with the specified owner
	cursor, err := listCollection.Find(ctx, bson.M{"owner": owner})
	if err != nil {
		return nil, fmt.Errorf("failed to find lists for owner '%s': %v", owner, err)
	}
	defer cursor.Close(ctx)

	var lists []string
	for cursor.Next(ctx) {
		var listDoc struct {
			List string `bson:"list"`
		}
		err := cursor.Decode(&listDoc)
		if err != nil {
			log.Printf("failed to decode list document: %v", err)
			continue
		}

		log.Printf("Checking list: %s", listDoc.List) // Log which list is being checked

		// Check if this list is associated with any of the specified targets
		found := false
		for _, target := range targets {
			log.Printf("Checking target: %s", target) // Log which target is being checked
			var targetDoc struct {
				List   string `bson:"list"`
				Target string `bson:"target"`
			}
			err = db.Collection("targets").FindOne(ctx, bson.M{"list": listDoc.List, "target": target}).Decode(&targetDoc)
			if err == nil { // If no error, the target is associated with the list
				lists = append(lists, listDoc.List)
				found = true
				log.Printf("Target matched: %s for list: %s", target, listDoc.List)
				break // No need to check other targets for this list
			} else {
				log.Printf("No match for target: %s in list: %s, error: %v", target, listDoc.List, err)
			}
		}

		if !found {
			log.Printf("No targets matched for list: %s", listDoc.List)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during cursor iteration: %v", err)
	}

	if len(lists) == 0 {
		return nil, fmt.Errorf("There are no lists matching the owner and at least one of the targets")
	}

	return lists, nil
}

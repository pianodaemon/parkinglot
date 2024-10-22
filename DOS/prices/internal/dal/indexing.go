package dal

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoUniqueIndex(db *mongo.Database, collName string, fieldName string) (string, error) {

	coll := db.Collection(collName)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{fieldName, 1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return "", err
	}

	return indexName, nil
}

func IndexCreationSamples(db *mongo.Database) error {

	indexName, err := CreateMongoUniqueIndex(db, "price_lists", "list")
	if err != nil {
		return err
	}
	fmt.Println("Index created for price_lists collection: ", indexName)

	indexName, err = CreateMongoUniqueIndex(db, "prices", "hash")
	if err != nil {
		return err
	}
	fmt.Println("Index created for prices collection: ", indexName)

	return nil
}

func CheckIndexesToCreate(db *mongo.Database) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// if there is no index for field "list", ascending, in collection "price_lists", then
	// a new pricing db is assumed and all indexes have to be created.
	collection := db.Collection("price_lists")

	// Check if the index already exists
	indexView := collection.Indexes()
	cursor, err := indexView.List(ctx)
	if err != nil {
		return err
	}

	var existingIndexes []bson.M
	if err = cursor.All(ctx, &existingIndexes); err != nil {
		return err
	}

	indexExists := false
	for _, index := range existingIndexes {

		if index["name"] == "list_1" {
			indexExists = true
			break
		}
	}

	if !indexExists {
		if err := IndexCreationSamples(db); err != nil {
			return err
		}
	} else {
		fmt.Println("Indexes already exist, skipping creation.")
	}

	return nil
}

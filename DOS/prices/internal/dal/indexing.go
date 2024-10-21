package dal

import (
	"context"
	"fmt"

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

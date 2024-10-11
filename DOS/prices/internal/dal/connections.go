package dal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Sets a mongo connection up
func SetUpConnMongoDB(mcli **mongo.Client, uri string) error {

	// Create a new MongoDB client
	cli, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// Set up a timeout for the connection
	ctxConn, cancelConn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConn()

	// Connect to MongoDB
	if err = cli.Connect(ctxConn); err != nil {
		return err
	}

	// Set up a timeout for the ping
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()

	// Ping the MongoDB server to ensure the connection is established
	if err = cli.Ping(ctxPing, readpref.Primary()); err != nil {
		return err
	}

	// Set the client if the connection and ping succeed
	*mcli = cli

	return nil
}

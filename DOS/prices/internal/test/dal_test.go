package verification

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"

	"blaucorp.com/prices/internal/dal"
)

func TestWithMongoDBContainer(t *testing.T) {
	// Request for a MongoDB container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0", // or the version you need
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(30 * time.Second),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer mongoC.Terminate(ctx)

	// Get the host and port for the MongoDB container
	host, err := mongoC.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %s", err)
	}

	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		t.Fatalf("failed to get container port: %s", err)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	var client *mongo.Client
	err = dal.SetUpConnMongoDB(&client, mongoURI)
	if err != nil {
		t.Fatalf("failed to set up mongo client: %s", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("pricing_db")

	//
	// Populate the data
	err = dal.CreatePriceList(db, "winter-2024-1728533139", "viajes Ponchito")
	if err != nil {
		panic(err.Error())
	}

	dal.AssignTargets(db, "winter-2024-1728533139", []string{"pepsi", "coca"})

	// Add fake prices
	dal.AddPrice(db, "winter-2024-1728533139", "1254-545-66", "m3", "madera", "limpia", 15.50)
	dal.AddPrice(db, "winter-2024-1728533139", "7845-155-78", "kg", "hierro", "sucia", 25.75)
	dal.AddPrice(db, "winter-2024-1728533139", "9987-845-23", "lt", "agua", "purificada", 3.80)

	// Retrieve a price by tuple
	priceTuple := map[string]string{
		"list":      "winter-2024-1728533139",
		"sku":       "1254-545-66",
		"unit":      "m3",
		"material":  "madera",
		"tservicio": "limpia",
	}

	price, err := dal.RetrievePriceByTuple(db, priceTuple)
	if err != nil {
		panic(err.Error())

	}

	if price != 15.50 {

		t.Fatalf("Price %.2f is not the expected", price)
	}
	//
	log.Println("MongoDB Testcontainer is running and connected")
}

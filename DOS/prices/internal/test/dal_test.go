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

	// Test cases for price retrieval
	priceTests := []struct {
		priceTuple    map[string]string
		expectedPrice float64
	}{
		{
			priceTuple: map[string]string{
				"list":      "winter-2024-1728533139",
				"sku":       "1254-545-66",
				"unit":      "m3",
				"material":  "madera",
				"tservicio": "limpia",
			},
			expectedPrice: 15.50,
		},
		{
			priceTuple: map[string]string{
				"list":      "winter-2024-1728533139",
				"sku":       "7845-155-78",
				"unit":      "kg",
				"material":  "hierro",
				"tservicio": "sucia",
			},
			expectedPrice: 25.75,
		},
		{
			priceTuple: map[string]string{
				"list":      "winter-2024-1728533139",
				"sku":       "9987-845-23",
				"unit":      "lt",
				"material":  "agua",
				"tservicio": "purificada",
			},
			expectedPrice: 3.80,
		},
	}

	// Loop through test cases to verify each price
	for _, test := range priceTests {
		price, err := dal.RetrievePriceByTuple(db, test.priceTuple)
		if err != nil {
			t.Fatalf("Failed to retrieve price for tuple %+v: %s", test.priceTuple, err)
		}

		if price != test.expectedPrice {
			t.Fatalf("Price %.2f for tuple %+v is not the expected %.2f", price, test.priceTuple, test.expectedPrice)
		}
	}

	log.Println("MongoDB Testcontainer is running and prices are verified")
}

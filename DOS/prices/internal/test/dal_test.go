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
	"blaucorp.com/prices/internal/misc"
)

func TestWithMongoDBContainer(t *testing.T) {

	ctx := context.Background()

	mongoC, err := setupMongoContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer mongoC.Terminate(ctx)

	mongoURI, err := getMongoURI(ctx, mongoC)
	if err != nil {
		t.Fatalf("Failed to retrieve MongoDB URI: %s", err)
	}

	var client *mongo.Client
	err = dal.SetUpConnMongoDB(&client, mongoURI)
	if err != nil {
		t.Fatalf("failed to set up mongo client: %s", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("pricing_db")

	listName := misc.GenerateNameWithCurrency("MXN", misc.GenerateNameWithTimestamp("winter"))
	verifyPrices(t, db, listName)
	verifyDeletion(t, db, listName)
}

// Helper function to set up MongoDB container
func setupMongoContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(30 * time.Second),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

// Helper function to get MongoDB URI
func getMongoURI(ctx context.Context, mongoC testcontainers.Container) (string, error) {
	host, err := mongoC.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		return "", fmt.Errorf("failed to get container port: %w", err)
	}

	return fmt.Sprintf("mongodb://%s:%s", host, port.Port()), nil
}

// Helper function to verify prices in the database
func verifyPrices(t *testing.T, db *mongo.Database, listName string) {

	owner := "viajes Ponchito"
	err := dal.CreatePriceList(db, listName, owner)
	if err != nil {
		t.Fatalf("Failed to create price list: %s", err)
	}

	targets := []string{"pepsi", "coca"}
	dal.AssignTargets(db, listName, targets)
	log.Printf("Assigned targets %v to list '%s'\n", targets, listName)

	availables, err := dal.GetListsByOwnerAndTargets(db, owner, targets)
	if err == nil {
		const uniqueSlotExpected = 0
		if len(availables) > 0 && availables[uniqueSlotExpected] != listName {
			t.Fatalf("The list was just created but it is not found %s", listName)
		}
	} else {
		t.Fatalf("List not retrieved: %s", err)
	}

	// Adding prices
	prices := []struct {
		sku, unit, material, tservicio string
		price                          float64
	}{
		{"1254-545-66", "m3", "madera", "limpia", 15.50},
		{"7845-155-78", "kg", "hierro", "sucia", 25.75},
		{"9987-845-23", "lt", "agua", "purificada", 3.80},
	}

	for _, p := range prices {
		err = dal.AddPrice(db, listName, p.sku, p.unit, p.material, p.tservicio, p.price)
		if err != nil {
			t.Fatalf("Failed to add price %v: %s", p, err)
		}
	}

	// Verify initial prices
	testCases := []struct {
		priceTuple    map[string]string
		expectedPrice float64
	}{
		{
			priceTuple: map[string]string{
				"list":      listName,
				"sku":       "1254-545-66",
				"unit":      "m3",
				"material":  "madera",
				"tservicio": "limpia",
			},
			expectedPrice: 15.50,
		},
		{
			priceTuple: map[string]string{
				"list":      listName,
				"sku":       "7845-155-78",
				"unit":      "kg",
				"material":  "hierro",
				"tservicio": "sucia",
			},
			expectedPrice: 25.75,
		},
		{
			priceTuple: map[string]string{
				"list":      listName,
				"sku":       "9987-845-23",
				"unit":      "lt",
				"material":  "agua",
				"tservicio": "purificada",
			},
			expectedPrice: 3.80,
		},
	}

	for _, test := range testCases {
		price, err := dal.RetrievePriceByTuple(db, test.priceTuple)
		if err != nil {
			t.Fatalf("Failed to retrieve price for tuple %+v: %s", test.priceTuple, err)
		}

		if price != test.expectedPrice {
			t.Fatalf("Expected price %.2f for tuple %+v, got %.2f", test.expectedPrice, test.priceTuple, price)
		}
	}

	log.Println("Initial price verification completed successfully")

	// **Price Edition Section**
	// Define updated prices
	updatedPrices := []struct {
		sku, unit, material, tservicio string
		newPrice                       float64
	}{
		{"1254-545-66", "m3", "madera", "limpia", 16.00},  // Updated from 15.50 to 16.00
		{"7845-155-78", "kg", "hierro", "sucia", 26.50},   // Updated from 25.75 to 26.50
		{"9987-845-23", "lt", "agua", "purificada", 4.00}, // Updated from 3.80 to 4.00
	}

	// Update prices
	for _, p := range updatedPrices {
		err = dal.EditPrice(db, listName, p.sku, p.unit, p.material, p.tservicio, p.newPrice)
		if err != nil {
			t.Fatalf("Failed to update price %v: %s", p, err)
		}
	}

	// Verify updated prices
	for _, test := range updatedPrices {
		priceTuple := map[string]string{
			"list":      listName,
			"sku":       test.sku,
			"unit":      test.unit,
			"material":  test.material,
			"tservicio": test.tservicio,
		}
		retrievedPrice, err := dal.RetrievePriceByTuple(db, priceTuple)
		if err != nil {
			t.Fatalf("Failed to retrieve updated price for tuple %+v: %s", priceTuple, err)
		}

		if retrievedPrice != test.newPrice {
			t.Fatalf("Expected updated price %.2f for tuple %+v, got %.2f", test.newPrice, priceTuple, retrievedPrice)
		}

		log.Printf("Price updated successfully for tuple %+v: %.2f\n", priceTuple, retrievedPrice)
	}

	log.Println("Price edition verification completed successfully")
}

// Helper function to verify deletion in the database
func verifyDeletion(t *testing.T, db *mongo.Database, listName string) {

	err := dal.DeleteList(db, listName)
	if err != nil {
		t.Fatalf("Failed to delete price list %s: %s", listName, err)
	}
}

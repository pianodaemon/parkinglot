package hookups

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"blaucorp.com/prices/internal/dal"
)

func Calis() {
	mongoURI := fmt.Sprintf("mongodb://user:123qwe@%s:%s/", "localhost", "27017")
	var client *mongo.Client
	err := dal.SetUpConnMongoDB(&client, mongoURI)
	if err != nil {
		panic("failed to set up mongo client: %s")
	}

	db := client.Database("pricing_db")

	// Populate the data
	err = dal.CreatePriceList(db, "invierno-2024-1728533139", "viajes Ponchito")
	if err != nil {
		panic(err.Error())
	}
	dal.AssignTargets(db, "invierno-2024-1728533139", []string{"pepsi", "coca"})

	// Add fake prices
	dal.AddPrice(db, "invierno-2024-1728533139", "1254-545-66", "m3", "madera", "limpia", 15.50)
	dal.AddPrice(db, "invierno-2024-1728533139", "7845-155-78", "kg", "hierro", "sucia", 25.75)
	dal.AddPrice(db, "invierno-2024-1728533139", "9987-845-23", "lt", "agua", "purificada", 3.80)

	// Retrieve a price by tuple
	priceTuple := map[string]string{
		"list":      "invierno-2024-1728533139",
		"sku":       "1254-545-66",
		"unit":      "m3",
		"material":  "madera",
		"tservicio": "limpia",
	}
	dal.RetrievePriceByTuple(db, priceTuple)

	ctx, cancelDisconn := context.WithTimeout(context.Background(), 2*time.Second)
	client.Disconnect(ctx)

	/* It'll be even called if succeeded just to
	   release resources of timing */
	defer cancelDisconn()
}

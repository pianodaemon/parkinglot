package main

import (
	"blaucorp.com/price-lists/dal"
)

func main() {
	// Connect to MongoDB
	client, ctx := dal.ConnectMongoDB()
	defer client.Disconnect(ctx)

	db := client.Database("pricing_db")

	// Populate the data
	dal.CreatePriceList(db, "invierno-2024-1728533139", "viajes Ponchito")
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
}

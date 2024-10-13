package hookups

import (
	"fmt"

	"blaucorp.com/prices/internal/dal"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	PricesManager struct {
		mcli *mongo.Client
	}
)

func NewPricesManager() *PricesManager {

	pm := &PricesManager{}

	// Connect to MongoDB along with pool of connections
	mongoURI := fmt.Sprintf("mongodb://user:123qwe@%s:%s/", "localhost", "27017")
	err := dal.SetUpConnMongoDB(&(pm.mcli), mongoURI)
	if err != nil {
		panic(err.Error())
	}

	return pm
}

func (self *PricesManager) DoCreatePriceList(listName, owner string) error {
	// Pending implementation
	return nil
}

func (self *PricesManager) DoDeleteList(listName string) error {
	// Pending implementation
	return nil
}

func (self *PricesManager) DoAssignTargets(listName string, targets []string) error {
	// Pending implementation
	return nil
}

func (self *PricesManager) DoUpdatePrice(listName, sku, unit, material, tservicio string, price float64) error {
	// Pending implementation
	return nil
}

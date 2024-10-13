package hookups

import (
	"fmt"

	"blaucorp.com/prices/internal/dal"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	PricesManager struct {
		dbID string
		mcli *mongo.Client
	}
)

func NewPricesManager() *PricesManager {

	pm := &PricesManager{}

	pm.dbID = "pricing_db"

	// Connect to MongoDB along with pool of connections
	mongoURI := fmt.Sprintf("mongodb://user:123qwe@%s:%s/", "localhost", "27017")
	err := dal.SetUpConnMongoDB(&(pm.mcli), mongoURI)
	if err != nil {
		panic(err.Error())
	}

	return pm
}

func (self *PricesManager) DoCreatePriceList(listName, owner string) error {
	db := self.mcli.Database(self.dbID)
	return dal.CreatePriceList(db, listName, owner)
}

func (self *PricesManager) DoDeleteList(listName string) error {
	db := self.mcli.Database(self.dbID)
	return dal.DeleteList(db, listName)
}

func (self *PricesManager) DoAssignTargets(listName string, targets []string) error {
	db := self.mcli.Database(self.dbID)
	return dal.AssignTargets(db, listName, targets)
}

func (self *PricesManager) DoEditPrice(listName, sku, unit, material, tservicio string, price float64) error {
	db := self.mcli.Database(self.dbID)
	return dal.EditPrice(db, listName, sku, unit, material, tservicio, price)
}

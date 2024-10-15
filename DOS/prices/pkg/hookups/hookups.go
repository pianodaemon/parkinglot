package hookups

import (
	"blaucorp.com/prices/internal/dal"

	"go.mongodb.org/mongo-driver/mongo"
)

type (

	// PricesManagerInterface defines the contract for managing price lists
	PricesManagerInterface interface {
		DoCreatePriceList(listName, owner string) error
		DoDeleteList(listName string) error
		DoAssignTargets(listName string, targets []string) error
		DoAddPrice(listName, sku, unit, material, tservicio string, price float64) error
		DoEditPrice(listName, sku, unit, material, tservicio string, price float64) error
		DoRetrievePriceByTuple(priceTuple map[string]string) (float64, error)
		DoGetListsByOwnerAndTargets(owner string, targets []string) ([]string, error)
	}

	PricesManager struct {
		dbID string
		mcli *mongo.Client
	}
)

func NewPricesManager(mongoURI string) *PricesManager {

	pm := &PricesManager{}

	pm.dbID = "pricing_db"
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

func (self *PricesManager) DoAddPrice(listName, sku, unit, material, tservicio string, price float64) error {
	db := self.mcli.Database(self.dbID)
	return dal.AddPrice(db, listName, sku, unit, material, tservicio, price)
}

func (self *PricesManager) DoEditPrice(listName, sku, unit, material, tservicio string, price float64) error {
	db := self.mcli.Database(self.dbID)
	return dal.EditPrice(db, listName, sku, unit, material, tservicio, price)
}

func (self *PricesManager) DoRetrievePriceByTuple(priceTuple map[string]string) (float64, error) {
	db := self.mcli.Database(self.dbID)
	return dal.RetrievePriceByTuple(db, priceTuple)
}

func (self *PricesManager) DoGetListsByOwnerAndTargets(owner string, targets []string) ([]string, error) {
	db := self.mcli.Database(self.dbID)
	return dal.GetListsByOwnerAndTargets(db, owner, targets)
}

package hookups

import ()

type (
	PricesManager struct{}
)

func NewPricesManager() *PricesManager {
	return &PricesManager{}
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

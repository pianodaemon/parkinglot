package hookups

import ()

type (
	PricesManager struct{}
)

func NewPricesManager() *PricesManager {
	return &PricesManager{}
}

func (self *PricesManager) DoUpdatePrice(listName, sku, unit, material, tservicio string, price float64) error {
	// Pending implementation
	return nil
}

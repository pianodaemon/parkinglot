package hookups

import ()

type (
	PricesManager struct{}
)

func NewPricesManager() *PricesManager {
	return &PricesManager{}
}

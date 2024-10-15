package service

import (
	co "blaucorp.com/prices/internal/controllers"
	hups "blaucorp.com/prices/pkg/hookups"
	"os"

	"github.com/gin-gonic/gin"
)

// Engages the RESTful API
func Engage() (merr error) {

	defer func() {

		if r := recover(); r != nil {
			merr = r.(error)
		}
	}()

	getEnv := func(key, fallback string) string {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
		return fallback
	}

	pricesManagerImplt := hups.NewPricesManager(getEnv("MONGO_URI", "mongodb://127.0.0.1:27017"))

	r := gin.Default()
	setUpHandlers(r, pricesManagerImplt)
	r.Run() // listen and serve on 0.0.0.0:8080

	return nil
}

func setUpHandlers(r *gin.Engine, pm *hups.PricesManager) {

	r.GET("/ping", co.Health)
	r.POST("/price-lists", co.CreateList(pm))
	r.PUT("/prices", co.UpdatePrice(pm))
	r.GET("/prices", co.RetrievePriceByTuple(pm))
	r.POST("/prices", co.AddPriceToList(pm))
	r.GET("/price-lists", co.GetListsByOwnerAndTargets(pm))
}

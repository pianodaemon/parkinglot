package service

import (
	"net/http"

	co "blaucorp.com/prices/internal/controllers"
	hups "blaucorp.com/prices/pkg/hookups"

	"github.com/gin-gonic/gin"
)

// Engages the RESTful API
func Engage() {

	pricesManagerImplt := hups.NewPricesManager()

	r := gin.Default()
	setUpHandlers(r, pricesManagerImplt)
	r.Run() // listen and serve on 0.0.0.0:8080
}

func setUpHandlers(r *gin.Engine, pm *hups.PricesManager) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/price-lists", co.CreateList(pm.DoCreatePriceList, pm.DoAssignTargets, pm.DoEditPrice))
	r.PUT("/prices", co.UpdatePrice)
}

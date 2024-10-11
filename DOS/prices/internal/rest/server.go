package rest

import (
	"fmt"
	"net/http"

	"blaucorp.com/prices/internal/dal"
	"github.com/gin-gonic/gin"
)

type price struct {
	Sku       string  `json:"sku" binding:"required"`
	Unit      string  `json:"unit" binding:"required"`
	Material  string  `json:"material" binding:"required"`
	Tservicio string  `json:"tservicio" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
}

type priceList struct {
	List    string   `json:"list" binding:"required"`
	Owner   string   `json:"owner" binding:"required"`
	Targets []string `json:"targets" binding:"required"`
	Prices  []price  `json:"prices" binding:"required"`
}

func GetEngine() *gin.Engine {
	return gin.Default()
}

func SetHandlers(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/price-lists", func(c *gin.Context) {
		reqPriceList := priceList{}

		if errP := c.ShouldBind(&reqPriceList); errP != nil {
			c.String(http.StatusBadRequest, "the body should be form of priceList type")
			return
		}
		// Connect to MongoDB
		client, ctx := dal.ConnectMongoDB()
		defer client.Disconnect(ctx)

		db := client.Database("pricing_db")

		// Populate the data
		dal.CreatePriceList(db, reqPriceList.List, reqPriceList.Owner)
		dal.AssignTargets(db, reqPriceList.List, reqPriceList.Targets)

		// Add fake prices
		for _, price := range reqPriceList.Prices {
			dal.AddPrice(db, reqPriceList.List, price.Sku, price.Unit, price.Material, price.Tservicio, price.Price)
		}

		fmt.Println("-------> body: ", reqPriceList)
		c.JSON(http.StatusOK, gin.H{
			"results": "ok",
		})
	})
}

package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	co "blaucorp.com/prices/internal/controllers"
	"blaucorp.com/prices/internal/dal"
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

	r.POST("/price-lists", co.CreateList(pm.DoCreatePriceList, pm.DoAssignTargets, pm.DoUpdatePrice))

	r.PUT("/prices", func(c *gin.Context) {

		type hashedPrice struct {
			Hash  string  `json:"hash" binding:"required"`
			Price float64 `json:"price" binding:"required"`
		}

		reqHashedPrice := hashedPrice{}

		if errP := c.ShouldBind(&reqHashedPrice); errP != nil {
			c.String(http.StatusBadRequest, "the body should be form of hashedPrice type")
			return
		}
		// Connect to MongoDB
		mongoURI := fmt.Sprintf("mongodb://user:123qwe@%s:%s/", "localhost", "27017")
		var client *mongo.Client
		err := dal.SetUpConnMongoDB(&client, mongoURI)
		if err != nil {
			panic(err.Error())
		}

		db := client.Database("pricing_db")

		err = dal.EditPrice(db, reqHashedPrice.Hash, reqHashedPrice.Price)
		if err != nil {
			panic(err.Error())
		}

		ctx, cancelDisconn := context.WithTimeout(context.Background(), 2*time.Second)
		client.Disconnect(ctx)

		/* It'll be even called if succeeded just to
		   release resources of timing */
		defer cancelDisconn()

		fmt.Println("-------> body: ", reqHashedPrice)
		c.JSON(http.StatusOK, gin.H{
			"results": "ok",
		})
	})

	r.PUT("/update", co.UpdatePrice)
}

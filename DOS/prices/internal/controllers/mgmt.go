package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"blaucorp.com/prices/internal/dal"
	"github.com/gin-gonic/gin"
)

type (
	price struct {
		Sku       string  `json:"sku" binding:"required"`
		Unit      string  `json:"unit" binding:"required"`
		Material  string  `json:"material" binding:"required"`
		Tservicio string  `json:"tservicio" binding:"required"`
		Price     float64 `json:"price" binding:"required"`
	}

	priceList struct {
		List    string   `json:"list" binding:"required"`
		Owner   string   `json:"owner" binding:"required"`
		Targets []string `json:"targets" binding:"required"`
		Prices  []price  `json:"prices" binding:"required"`
	}

	priceUpdateRequest struct {
		price        // Embedding price to inherit its fields
		List  string `json:"list" binding:"required"`
	}

	DoAssignTargetsHandler   func(listName string, targets []string) error
	DoCreatePriceListHandler func(listName, owner string) error
	DoUpdatePriceHandler     func(listName, sku, unit, material, tservicio string, price float64) error
)

func CreateList(dclh DoCreatePriceListHandler, dath DoAssignTargetsHandler, duph DoUpdatePriceHandler) func(c *gin.Context) {

	return func(c *gin.Context) {

		reqPriceList := priceList{}

		if errP := c.ShouldBind(&reqPriceList); errP != nil {
			c.String(http.StatusBadRequest, "the body should be form of priceList type")
			return
		}

		err := dclh(reqPriceList.List, reqPriceList.Owner)
		if err != nil {
			panic(err.Error())
		}
		dath(reqPriceList.List, reqPriceList.Targets)

		// Add fake prices
		for _, price := range reqPriceList.Prices {
			duph(reqPriceList.List, price.Sku, price.Unit, price.Material, price.Tservicio, price.Price)
		}

		fmt.Println("-------> body: ", reqPriceList)
		c.JSON(http.StatusOK, gin.H{
			"results": "ok",
		})
	}
}

func UpdatePrice(c *gin.Context) {
	var req priceUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Connect to MongoDB
	mongoURI := fmt.Sprintf("mongodb://user:123qwe@%s:%s/", "localhost", "27017")
	var client *mongo.Client
	err := dal.SetUpConnMongoDB(&client, mongoURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	db := client.Database("pricing_db")

	// Edit the price using dal.AddOrUpdatePrice
	err = dal.AddOrUpdatePrice(db, req.List, req.Sku, req.Unit, req.Material, req.Tservicio, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update price"})
		return
	}

	// Disconnect MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client.Disconnect(ctx)

	c.JSON(http.StatusOK, gin.H{"message": "Price updated successfully"})
}

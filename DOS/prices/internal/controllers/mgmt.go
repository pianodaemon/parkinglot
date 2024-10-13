package controllers

import (
	"fmt"
	"net/http"

	hups "blaucorp.com/prices/pkg/hookups"

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
)

func CreateList(pricesManagerImplt hups.PricesManagerInterface) func(c *gin.Context) {

	return func(c *gin.Context) {

		reqPriceList := priceList{}

		if errP := c.ShouldBind(&reqPriceList); errP != nil {
			c.String(http.StatusBadRequest, "the body should be form of priceList type")
			return
		}

		err := pricesManagerImplt.DoCreatePriceList(reqPriceList.List, reqPriceList.Owner)
		if err != nil {
			panic(err.Error())
		}
		pricesManagerImplt.DoAssignTargets(reqPriceList.List, reqPriceList.Targets)

		// Add fake prices
		for _, price := range reqPriceList.Prices {
			pricesManagerImplt.DoAddPrice(reqPriceList.List, price.Sku, price.Unit, price.Material, price.Tservicio, price.Price)
		}

		fmt.Println("-------> body: ", reqPriceList)
		c.JSON(http.StatusOK, gin.H{
			"results": "ok",
		})
	}
}

func UpdatePrice(pricesManagerImplt hups.PricesManagerInterface) func(c *gin.Context) {

	return func(c *gin.Context) {
		var req priceUpdateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		if err := pricesManagerImplt.DoEditPrice(req.List, req.Sku, req.Unit, req.Material, req.Tservicio, req.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update price"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Price updated successfully"})
	}
}

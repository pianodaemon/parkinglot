package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	// PricesManagerInterface defines the contract for managing price lists
	PricesManagerInterface interface {
		DoCreatePriceList(listName, owner, currency string) (string, error)
		DoDeleteList(listName string) error
		DoAssignTargets(listName string, targets []string) error
		DoAddPrice(listName, sku, unit, material, tservicio string, price float64) error
		DoEditPrice(listName, sku, unit, material, tservicio string, price float64) error
		DoRetrievePriceByTuple(priceTuple map[string]string) (float64, error)
		DoGetListsByOwnerAndTargets(owner string, targets []string) ([]string, error)
	}

	price struct {
		Sku       string  `json:"sku" binding:"required"`
		Unit      string  `json:"unit" binding:"required"`
		Material  string  `json:"material" binding:"required"`
		Tservicio string  `json:"tservicio" binding:"required"`
		Price     float64 `json:"price" binding:"required"`
	}

	priceList struct {
		List     string   `json:"list" binding:"required"`
		Owner    string   `json:"owner" binding:"required"`
		Currency string   `json:"currency" binding:"required"`
		Targets  []string `json:"targets" binding:"required"`
		Prices   []price  `json:"prices" binding:"required"`
	}

	priceUpdateRequest struct {
		price        // Embedding price to inherit its fields
		List  string `json:"list" binding:"required"`
	}
)

func CreateList(pricesManagerImplt PricesManagerInterface) func(c *gin.Context) {

	return func(c *gin.Context) {

		reqPriceList := priceList{}

		if errP := c.ShouldBind(&reqPriceList); errP != nil {
			c.String(http.StatusBadRequest, "the body should be form of priceList type")
			return
		}

		name, err := pricesManagerImplt.DoCreatePriceList(reqPriceList.List, reqPriceList.Owner, reqPriceList.Currency)
		if err != nil {
			panic(err.Error())
		}
		pricesManagerImplt.DoAssignTargets(name, reqPriceList.Targets)

		// Add fake prices
		for _, price := range reqPriceList.Prices {
			pricesManagerImplt.DoAddPrice(name, price.Sku, price.Unit, price.Material, price.Tservicio, price.Price)
		}

		fmt.Println("-------> body: ", name)
		c.JSON(http.StatusOK, gin.H{
			"results": "ok",
		})
	}
}

func UpdatePrice(pricesManagerImplt PricesManagerInterface) func(c *gin.Context) {

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

func RetrievePriceByTuple(pricesManagerImplt PricesManagerInterface) func(c *gin.Context) {

	return func(c *gin.Context) {
		list := c.Query("list")
		sku := c.Query("sku")
		unit := c.Query("unit")
		material := c.Query("material")
		tservicio := c.Query("tservicio")

		priceTupl := map[string]string{
			"list":      list,
			"sku":       sku,
			"unit":      unit,
			"material":  material,
			"tservicio": tservicio,
		}

		priceVal, err := pricesManagerImplt.DoRetrievePriceByTuple(priceTupl)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"precio": fmt.Sprintf("%.2f", priceVal)})
	}
}

func AddPriceToList(pricesManagerImplt PricesManagerInterface) func(c *gin.Context) {

	return func(c *gin.Context) {
		var req priceUpdateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		if err := pricesManagerImplt.DoAddPrice(req.List, req.Sku, req.Unit, req.Material, req.Tservicio, req.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add price"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Price added successfully"})
	}
}

func GetListsByOwnerAndTargets(pricesManagerImplt PricesManagerInterface) func(c *gin.Context) {

	return func(c *gin.Context) {
		reqOwner := c.Query("owner")
		reqTargets := c.QueryArray("targets")

		lists, err := pricesManagerImplt.DoGetListsByOwnerAndTargets(reqOwner, reqTargets)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"lists": lists})
	}
}

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	hups "blaucorp.com/fiscal-engine/pkg/hookups"
)

type (
	FiscalEngineInterface interface {
		DoCreateReceipt(dto *hups.ReceiptDTO) (string, error)
	}
)

func CreateReceipt(fiscalEngineImplt FiscalEngineInterface) func(c *gin.Context) {

	return func(c *gin.Context) {

		dto := &hups.ReceiptDTO{}

		if errP := c.ShouldBind(dto); errP != nil {
			c.String(http.StatusBadRequest, "the body should be form of recipe type")
			return
		}

		name, err := fiscalEngineImplt.DoCreateReceipt(dto)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("-------> body: ", name)
		c.JSON(http.StatusOK, gin.H{
			"results": "ok",
		})
	}
}

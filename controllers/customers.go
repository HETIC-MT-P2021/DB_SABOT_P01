package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"packages.hetic.net/gosql/models"
)

// GetCustomer handle request to call order handler
func GetCustomer(c *gin.Context) {
	customerNumber := c.Param("customerNumber")

	var customer models.Customer

	customer, err := models.GetCustomer(customerNumber, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err,
			"data":    false,
		})
	} else {
		c.JSON(http.StatusFound, gin.H{
			"message": "Found customer successfully",
			"success": true,
			"data":    customer,
		})
	}
}

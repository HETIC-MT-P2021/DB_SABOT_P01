package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"packages.hetic.net/gosql/models"
)

// GetOrder handle request to call order handler
func GetOrder(c *gin.Context) {
	orderNumber := c.Param("orderNumber")

	var order models.Order

	order, err := models.GetOrderDetails(orderNumber, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err,
			"data":    false,
		})
	} else {
		c.JSON(http.StatusFound, gin.H{
			"message": "Found order successfully",
			"success": true,
			"data":    order,
		})
	}
}

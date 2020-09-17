package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"packages.hetic.net/gocqrs/domain"
	"packages.hetic.net/gocqrs/models"
)

// CreateOrder handle request to call order handler
func CreateOrder(c *gin.Context) {
	createdBy := c.PostForm("createdBy")

	var order models.Instruction

	order.Data = createdBy

	order.Operation = "Create"

	orderCreated := domain.OrderBusCommandHandler(order)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Created order successfully",
		"content": orderCreated,
	})
}

// GetOrders handle request to call order handler
func GetOrders(c *gin.Context) {
	var order models.Instruction

	order.Operation = "Selects"
	order.Data = c

	ordersInstruction := domain.OrderBusQueryHandler(order)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Found orders successfully",
		"content": ordersInstruction.Data.([]models.CreatedOrder),
	})
}

// GetOrder handle request to call order handler
func GetOrder(c *gin.Context) {
	orderID := c.Param("orderID")

	var order models.Instruction

	order.Operation = "Select"
	order.Data = orderID

	ordersInstruction := domain.OrderBusQueryHandler(order)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Found orders successfully",
		"content": ordersInstruction.Data.([]models.CreatedOrder),
	})
}

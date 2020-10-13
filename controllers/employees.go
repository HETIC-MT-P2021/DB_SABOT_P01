package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"packages.hetic.net/gosql/models"
)

// GetEmployees handle request to call order handler
func GetEmployees(c *gin.Context) {
	var employees []models.Employee

	employees, err := models.GetEmployees(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err,
			"data":    false,
		})
	} else {
		c.JSON(http.StatusFound, gin.H{
			"message": "Found employees successfully",
			"success": true,
			"data":    employees,
		})
	}
}

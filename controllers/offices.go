package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"packages.hetic.net/gosql/models"
)

// GetOffice handle request to call order handler
func GetOffice(c *gin.Context) {
	officeNumber := c.Param("officeNumber")

	var office models.Office

	office, employees, err := models.GetOffice(officeNumber, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err,
			"data":    false,
		})
	} else {
		c.JSON(http.StatusFound, gin.H{
			"message": "Found office successfully",
			"success": true,
			"data":    struct { Office models.Office; Employees []models.SimplifiedEmployee} { office, employees },
		})
	}
}

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"packages.hetic.net/gosql/controllers"
)

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "API is running successfully",
		"success": true,
	})
}

// StartRouter will launch the web server
func StartRouter(apiPort string) {
	router := gin.New()

	public := router.Group("/")
	{
		public.GET("/customer/:customerNumber", controllers.GetCustomer)

		public.GET("/orders/:orderNumber", controllers.GetOrder)

		public.GET("/employees", controllers.GetEmployees)

		public.GET("/offices/:officeNumber", controllers.GetOffice)
	}

	router.Run(":" + apiPort)
}

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"packages.hetic.net/gocqrs/controllers"
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
		public.POST("/orders", controllers.CreateOrder)

		public.GET("/orders", controllers.GetOrders)
		public.GET("/orders/:orderID", controllers.GetOrder)
	}

	router.Run(":" + apiPort)
}

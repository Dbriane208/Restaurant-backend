package routes

// importing the necessary packages
import (
	controller "restaurant-backend/controllers"

	"github.com/gin-gonic/gin"
)
// function responsible for configuring the routes related to order operations
// takes an argument,incomingRoutes of type *gin.Engine
func OrderRoutes(incomingRoutes *gin.Engine){
	// Get request that retrieves a list of orders from the database
	incomingRoutes.GET("/orders",controller.GetOrders())
	// Get request that retrieves a specific order from the database
	incomingRoutes.GET("/orders/:order_id",controller.GetOrder())
	// Post request that creates a new order to the database
	incomingRoutes.POST("/orders",controller.CreateOrder())
	// Patch request that updates a specific order entry from the database
	incomingRoutes.PATCH("/orders/:order_id",UpdateOrder())
}
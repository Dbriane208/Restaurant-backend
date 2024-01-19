package routes

// importing the necessary packages
import (
	controller "restaurant-backend/controllers"

	"github.com/gin-gonic/gin"
)

// function responsible for configuring routes related to orderItem operation
// takes an argument of type *gin.Engine
func OrderItemRoutes(incomingRoutes *gin.Engine){
	// Get request that retrieves a list of order items from the database
	incomingRoutes.GET("/orderItems",controller.GetOrderItems())
	// Get request that retrives a specific item from the database
	incomingRoutes.GET("/orderItems/:orderItem_id",controller.GetOrderItem())
	// Get request that retrieves a specific order from the database
	incomingRoutes.GET("/orderItems-order/:order_id",controller.GetOrderItemsByOrder())
	// Post request that creates a new order entry to the database
	incomingRoutes.POST("/orderItems",controller.CreateOrderItem())
	// Patch request that updates a specific order item entry
	incomingRoutes.PATCH("/orderItems/:orderItem_id",controller.UpdateOrderItem())
}
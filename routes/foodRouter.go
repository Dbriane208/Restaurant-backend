package routes

import (
	controller "restaurant-backend/controllers"

	"github.com/gin-gonic/gin"
)

// function responsible for configuring routes related to food operation
// It takes a gin engine argument, incomingRoutes.
func FoodRoutes(incomingRoutes *gin.Engine){
	// the Get request retrives a list of foods from the database
	incomingRoutes.GET("/foods",controller.GetFoods())
	// the Get request retrieves  a specific type of food from the database
	incomingRoutes.GET("/foods/:food_id",controller.GetFood())
	// the Post request creates a new food item in the database 
	incomingRoutes.POST("/foods",controller.CreateFood())
	// the Patch request updates a specific item in the database
	incomingRoutes.PATCH("/foods/:food_id",controller.UpdateFood())
}
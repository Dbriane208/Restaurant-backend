package routes

// importing the necessary libraries
import (
	controller "restaurant-backend/controllers"

	"github.com/gin-gonic/gin"
)

// function responsible for configuring routes related to Menu operations
// it takes an argument of type *gin.Engine
func MenuRoutes(incomingRoutes *gin.Engine){
	// Get request that retrieves a list of menus
	incomingRoutes.GET("/menus",controller.GetMenus())
	// Get request that retrieves a specific menu
	incomingRoutes.GET("/menus/:menu_id",controller.GetMenu())
	// Post request that creates a new menu into the database
	incomingRoutes.POST("/menus",controller.CreateMenu())
	// Patch request that updates a menus specific entry
	incomingRoutes.PATCH("/menus/:menu_id",controller.UpdateMenu())
}
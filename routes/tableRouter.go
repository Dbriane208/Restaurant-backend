package routes

// importing the necessary packages
import (
	controller "restaurant-backend/controllers"

	"github.com/gin-gonic/gin"
)

// function responsible for configuring routes related to tables operations
// takes an argument of type *gin.Engine
func TableRoutes(incomingRoutes *gin.Engine){
	// Get request that retreives a list of tables from the database
	incomingRoutes.GET("/tables",controller.GetTables())
	// Get request that retrieves a specific table from the database
	incomingRoutes.GET("/tables/:table_id",controller.GetTable())
	// Post request that creates a new table entry in the database
	incomingRoutes.POST("/tables",controller.CreateTable())
	// Patch request that updates a specific entry in the database
	incomingRoutes.PATCH("/tables/:table_id",controller.UpdateTable())
}
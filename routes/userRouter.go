package routes

import (
	controller "restaurant-backend/controllers"

	"github.com/gin-gonic/gin"
)

// function responsible for configuring the user operations
// the function takes a gin engine argument,incomingRoutes
func UserRoutes(incomingRoutes *gin.Engine){
	// the Get request retrieves a list of users from the database
	incomingRoutes.GET("/users",controller.GetUsers())
	// the Get request retrieves a specific user from the database
	incomingRoutes.GET("/users/:user_id",controller.GetUser())
	// the Post request creates a new user to the database
	incomingRoutes.POST("/users/signup",controller.SignUp())
	// the Post request creates the user to the database
	incomingRoutes.POST("/users/login",controller.Login())
}
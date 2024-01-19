package main

import (
	"os"
	"restaurant-backend/database"
	"restaurant-backend/middleware"
	"restaurant-backend/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// a mongoDB collection for handling food related data
// the open collection function establishes the connection to the database
var foodCollection *mongo.Collection = database.OpenCollection(database.Client,"food")

func main(){

	// retrieves the value of the Port environment variable. Default port is 8000
	// flexible as the app can run on the specified port
	port := os.Getenv("PORT")
	if port == ""{
       port = "8000"
	}

	// Initializes the gin router and adds a logging middleware to log
	// HTTP requests.
	router := gin.New()
	router.Use(gin.Logger())

	// configures routes related to user operations by calling routes
	routes.UserRoutes(router)
	// Adds authentication middleware to the router that checks if requests are properly authenicated
	router.Use(middleware.Authentication())

	// configures variables routes for various operations by calling corresponding functions
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.InvoiceRoutes(router)
	routes.OrderItemRoutes(router)

	// Starts the HTTP server and listens on the specified port
	// The application will now handle incoming HTTP requests based on the configured routes
	router.Run(":" + port)
	
}
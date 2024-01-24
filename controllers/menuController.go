package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-backend/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// gin.Handler/func represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

// opens a MongoDb collection specified by the name menu
var menuCollection *mongo.Collection = database.OpenCollection(database.Client,"menu")

func GetMenus() gin.HandlerFunc{
	// Handler function for getting the menu items
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		
		// Finding all the documents in the menuCollection.
		// bson.M{} is an empty filter indicating that all documents should be retrieved
		result,err :=  menuCollection.Find(context.TODO(),bson.M{})

		// Ensures that the sorrounding context is closed when the GetMenus completes
		defer cancel()

		// Checks if an error occured during the mongodb operation. If there's the error returns a json
		// and function exits
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing the menu item"})
		}

		// Declaring a slice to store the retrieved menu items
		var allMenus []bson.M

		// Extract all documents,from the mongodb result into the  "allmenu" slice
		if err = result.All(ctx,&allMenus); err != nil {
			log.Fatal(err)
		}

		// If everything is successful, the retrieved menu items are returned as a JSON response
        c.JSON(http.StatusOK,allMenus)
	}
}

func GetMenu() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func CreateMenu() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func UpdateMenu() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}
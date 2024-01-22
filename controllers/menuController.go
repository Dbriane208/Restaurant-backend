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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client,"menu")

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		result,err :=  menuCollection.Find(context.TODO(),bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing the menu item"})
		}

		var allMenus []bson.M
		if err = result.All(ctx,&allMenus); err != nil {
			log.Fatal(err)
		}
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
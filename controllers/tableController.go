package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-backend/database"
	"restaurant-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// gin.HandlerFunc represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

var tableCollection *mongo.Collection = database.OpenCollection(database.Client,"table")

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		result,err := tableCollection.Find(context.TODO(),bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing table items"})
		}

		var allTables []bson.M
		if err = result.All(ctx,&allTables); err != nil{
			log.Fatal(err)
		}

		c.JSON(http.StatusOK,allTables)
	}
}

func GetTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var table models.Table
		tableId := c.Param("table_id")

		err := tableCollection.FindOne(ctx,bson.M{"table_id":tableId}).Decode(&table)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error when getting a table"})
			return
		}

		c.JSON(http.StatusOK,table)

	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		
	}
}
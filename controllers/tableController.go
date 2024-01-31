package controllers

import (
	"restaurant-backend/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// gin.HandlerFunc represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

var tableCollection *mongo.Collection = database.OpenCollection(database.Client,"table")

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func GetTable() gin.HandlerFunc{
	return func(c *gin.Context) {

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
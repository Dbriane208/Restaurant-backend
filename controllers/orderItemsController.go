package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// gin.HandlerFunc represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

func GetOrderItems() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func GetOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func ItemByOrder(id string) (orderItems []primitive.M, err error){

} 

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}
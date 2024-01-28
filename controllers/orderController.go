package controllers

import "github.com/gin-gonic/gin"

// gin.HandlerFunc represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		
	}
}
package controllers

import "github.com/gin-gonic/gin"

// gin.HandlerFunc represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

func GetTables() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func GetTable() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func CreateTable() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func UpdateTable() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}
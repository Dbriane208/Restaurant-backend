package controllers

import "github.com/gin-gonic/gin"

// gin.Handler/func represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

func GetMenus() gin.HandlerFunc{
	return func(ctx *gin.Context) {

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
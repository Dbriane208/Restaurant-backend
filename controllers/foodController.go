package controllers

import "github.com/gin-gonic/gin"

// gin.handlerFunc represents a request handler in gin
// func(ctx *gin.Context) represents an anonymous function that actuallyu handles the request
// [ctx *gin.Context] represents the actual HTTP request and response

func GetFoods() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func GetFood() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func CreateFood() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func round(num float64)int{

}

func toFixed(num float64,person int)float64{

}

func UpdateFood() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}
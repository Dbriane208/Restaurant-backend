package controllers

import "github.com/gin-gonic/gin"
// gin.HandlerFunc represents a request handler in gin
// func(ctx *gin.Context) represents an anonymous function which handles the actual route request
// [ctx *gin.Context] represents the current HTTP request and response

func GetInvoices() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func GetInvoice() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}
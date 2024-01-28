package controllers

import "github.com/gin-gonic/gin"

// gin.HandlerFunc represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func HashPassword(password string) string{

}

func verifyPassword(userPassword string,providePassword string)(bool,string){

}
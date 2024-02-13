package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-backend/database"
	helper "restaurant-backend/helpers"
	"restaurant-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// gin.HandlerFunc represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

var userCollection *mongo.Collection = database.OpenCollection(database.Client,"users")

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		recordPerPage,rpgError := strconv.Atoi(c.Query("recordPerPage"))
		if rpgError != nil || recordPerPage < 1{
			recordPerPage = 10
		}

		page, pageError := strconv.Atoi(c.Query("page"))
		if pageError != nil || page < 1{
			page = 1
		}

		startIndex := (page-1)* recordPerPage
		startIndex,_ = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match",Value: bson.D{{}}}}
		projectStage := bson.D{
			{Key: "$project",Value: bson.D{
				{Key: "_id",Value: 0},
				{Key: "total_count",Value: 1},
				{Key: "user_items",Value: bson.D{{Key: "$slice",Value: []interface{}{"$data",startIndex,recordPerPage}}}},
			}},
		}

		result,err := userCollection.Aggregate(ctx,mongo.Pipeline{
			matchStage,
			projectStage,
		})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing user items"})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx,&allUsers); err != nil{
			log.Fatal(err)
		}
		
		c.JSON(http.StatusOK,allUsers[0])
	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context) {
        
	}
}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var user models.User

		// convert the JSON data coming from postman to something golang understands

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		// validate the data based on user struct

		validationErr := validate.Struct(user)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}

		// You'll check if the email has already been used by another user
  
		emailCount,err := userCollection.CountDocuments(ctx,bson.M{"email":user.Email})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"erro":"error occured while checking for the email"})
			return
		}

		if emailCount > 0{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"this email already exists"})
		}
		// hash password

		password := HashPassword(*user.Password)
		user.Password = &password

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}	// You'll also check if the phone number has already been used by another person

		phoneCount,err := userCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
		defer cancel()
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while checking for the phone number"})
			return
		}

		if phoneCount > 0{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"this phone number already exists"})
		}

		// Create some extra details for the user object - created_at,updated_at, ID
		user.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// Generate token and refresh token(generate all tokens function helper)
		token,refreshToken := helper.GenerateAllTokens(*user.Email,*user.First_name,*user.Last_name,*&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		// If all OK, then you insert this new user into the user collection
		resultInsertionNumber,insertErr := userCollection.InsertOne(ctx,user)
		if insertErr != nil{
			msg := fmt.Sprintf("user item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"email":msg})
			return
		}

		// returns status OK and send the result back
		c.JSON(http.StatusOK,resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {

		// convert the login data from postman which is in JSON to golang readable format

		// find a user with that email and see if that user even exists

		// then you will verify the password

		// if all goes well then you'll generate tokens

		// Update tokens - tokens and refresh token

		// return OK
	}
}

func HashPassword(password string) string{

}

func verifyPassword(userPassword string,providePassword string)(bool,string){

}
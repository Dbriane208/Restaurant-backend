package controllers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"restaurant-backend/database"
	"restaurant-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// gin.handlerFunc represents a request handler in gin
// func(ctx *gin.Context) represents an anonymous function that actually handles the request
// [ctx *gin.Context] represents the actual HTTP request and response

// function that creates a food collection in the database
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

// used to struct field validation
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		recordPerPage,err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1{
			recordPerPage = 10
		}

		page,err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1{
			page = 1
		}

		startIndex := (page-1)*recordPerPage
		startIndex,err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key:"$match",Value:bson.D{{}}}}
		groupStage := bson.D{
			{Key:"$group",Value:bson.D{{Key:"_id",Value:"null"}}},
			{Key:"total_count",Value:bson.D{{Key:"$sum",Value:1}}},
			{Key:"data",Value:bson.D{{Key:"$push",Value:"$$ROOT"}}},
		}
		projectStage := bson.D{
			{Key:"$project",Value: bson.D{
				{Key: "_id",Value: 0},
				{Key: "total_count",Value: 1},
				{Key: "food_items",Value: bson.D{{Key: "$slice",Value :[]interface{}{"$data",startIndex,recordPerPage}}}},
			}},
		}

		result,err := foodCollection.Aggregate(ctx,mongo.Pipeline{
			matchStage,groupStage,projectStage,
		})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occurred while listing food items"})
		}

		var allFoods []bson.M
		if err = result.All(ctx,&allFoods); err != nil{
			log.Fatal(err)
		}

		c.JSON(http.StatusOK,allFoods[0])
	}
}

func GetFood() gin.HandlerFunc {
	// function to handle the food item
	return func(c *gin.Context) {
		// Creates a context with a timeout of 100 seconds
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// Retrieves the value of the "food_id" parameter from the request
		foodId := c.Param("food_id")
		// Create an instance of the "Food" struct from the "models" package
		var food models.Food

		// Mongodb's "FindOne" method to query for a document with a matching "food_id"
		// in the "foodCollection". Decode the result into the "food" variable
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		// Ensures the sorrounding context is closed when the GetFood completes
		defer cancel()

		// Handles any errors that occurred during the MongoDB query,returns an internal server error incase
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})
		}

		// Error is nil so it  returns the fetched food item as a JSON response
		// with a status of OK
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	// function that handles the creation of food request
	return func(c *gin.Context) {
		// creating a context with a time out of 100 seconds
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// creating the instances of the food and menu structs
		var menu models.Menu
		var food models.Food

		// binding the Json data from the Http request body to the food variable
		// and returning an error if data is not nil
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// checking whether the data entered in the struct from the Http is in the right format
		// the validate functions returns json error with status bad request
		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Querying the menu collection to find a menu based on food.menu_id is there
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		// cancelling the context after the menu querying
		defer cancel()

		// Handling the error incase the querying is not successful
		if err != nil {
			msg := fmt.Sprint("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// creating the timestamps when the food is created
		// the parse function receives a layout and a value, time.RFC3339, time.Now().Format(time.RFC3339)
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// creating a new id for the food collection
		// the hex function assigns a hexadecimal representation to the generated id
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()

		// Rounding up the food price to two decimal places
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		// Inserting the food struct into the food collection
		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			// Error handling incase the insertion has failed
			msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
        // canceling the context after a insertion to the database
		defer cancel()

		// returning the result of the insertion as the response
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5,num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10,float64(precision))
	return float64(round(num*output)) / output
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var menu models.Menu
		var food models.Food

		foodId := c.Param("food_id")
		if err := c.BindJSON(&food); err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}

		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key :"name",Value: food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "price",Value: food.Price})
		}

		if food.Food_image != nil{
			updateObj = append(updateObj, bson.E{Key: "food_image",Value: food.Food_image})
		}

		if food.Menu_id != nil {
			err := menuCollection.FindOne(ctx,bson.M{"menu_id":food.Menu_id}).Decode(&menu)
			defer cancel()
			if err != nil{
				msg := fmt.Sprintf("message:Menu was not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				return
			}

			updateObj = append(updateObj, bson.E{Key: "menu",Value: food.Price})
		}

		food.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at",Value: food.Updated_at})

		upsert := true
		filter := bson.M{"food_id":foodId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result,err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set",Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprintf("food item update failed")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		c.JSON(http.StatusOK,result)

	}
}

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
	// Handles the actual request of the food items
	return func(c *gin.Context) {
		// creates a context with a time out of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// parsing the query parameter "recordPerPage" from the request and convert it to an integer
		// the Atoi function converts a string to an integer
		recordPerPage,err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1{
			// settting the record per page to default of 10 incase of invalid value
			recordPerPage = 10
		}

		// parses the query parameter "page" from the request and converts it to an integer
		page,err := strconv.Atoi(c.Query("page"))
		// sets the default page to 1 incase of an invalid value
		if err != nil || page < 1{
			page = 1
		}

		// calculating the startIndex for pagination
		// parses the query "startIndex" for the starting index of the paginated result
		startIndex := (page-1)*recordPerPage
		startIndex,err = strconv.Atoi(c.Query("startIndex"))

		// Defines a MongoDB aggregation pipeline stage to match all documents
		matchStage := bson.D{{Key:"$match",Value:bson.D{{}}}}
		// Defines a mongoDB aggregation pipeline stage for grouping calculating total count and 
		// pushing data for pagination
		groupStage := bson.D{
			{Key:"$group",Value:bson.D{{Key:"_id",Value:"null"}}},
			{Key:"total_count",Value:bson.D{{Key:"$sum",Value:1}}},
			{Key:"data",Value:bson.D{{Key:"$push",Value:"$$ROOT"}}},
		}
		// Defines a mongoDB aggregation pipeline stage for projecting the result including the total count and
		// pushing data for pagination
		projectStage := bson.D{
			{Key:"$project",Value: bson.D{
				{Key: "_id",Value: 0},
				{Key: "total_count",Value: 1},
				{Key: "food_items",Value: bson.D{{Key: "$slice",Value :[]interface{}{"$data",startIndex,recordPerPage}}}},
			}},
		}

		// performs the mongoDB aggregation using the defined pipeline stages
		result,err := foodCollection.Aggregate(ctx,mongo.Pipeline{
			matchStage,groupStage,projectStage,
		})

		// Defering the cancelation of the context until the function returns 
		defer cancel()

		// Returning an internal server error incase the aggregation has failed
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occurred while listing food items"})
		}

		// creates a variable to store all the food result
		var allFoods []bson.M
		// retriving all the documents from the aggregation result and logs any errors
		if err = result.All(ctx,&allFoods); err != nil{
			log.Fatal(err)
		}

		// responds with the paginated food items in JSON format
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

// function that rounds a floating-point number to the nearest integer
func round(num float64) int {
	return int(num + math.Copysign(0.5,num))
}
// function that fixes the number of decimal places for a floating-point number
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10,float64(precision))
	return float64(round(num*output)) / output
}

func UpdateFood() gin.HandlerFunc {
	// function that handles the actual update food request items
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// creating instances of menu and food structs
		var menu models.Menu
		var food models.Food

		// retrieving the "food_id" from the HTTP request
		foodId := c.Param("food_id")

		// extracting and decoding JSON data from the HTTP request body to the food struct
		if err := c.BindJSON(&food); err != nil{
			// returning an internal server error incase the extraction and decoding fails
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}

		// creating a variable to store updates operations in a BSON document
		var updateObj primitive.D

		// appending the name to the updateObj if it's not null
		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key :"name",Value: food.Name})
		}

		// appending the price to the updateObj if it's not null
		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "price",Value: food.Price})
		}

		// appending the Food image to the updateObj if it's not null
		if food.Food_image != nil{
			updateObj = append(updateObj, bson.E{Key: "food_image",Value: food.Food_image})
		}

		// appending the Menu id to the updateObj if it's not null
		if food.Menu_id != nil {
			// Quering the database to find the document with the correspoding id
			err := menuCollection.FindOne(ctx,bson.M{"menu_id":food.Menu_id}).Decode(&menu)
			// defering the cancelation of the context until it returns
			defer cancel()
			// returning an internal error server incase the query operation fails
			if err != nil{
				msg := fmt.Sprintf("message:Menu was not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				return
			}

			// updating the menu_id
			//TODO
			updateObj = append(updateObj, bson.E{Key: "menu_id",Value: food.Menu_id})
		}

		// Updating the Updated at time with the current time
		food.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at",Value: food.Updated_at})

		// upsert is the combination of update and insert- used when you 
		// want to insert or update the document
		upsert := true
		// creating a mongoDb filter by creating a map where the key is
		//  "food_id" and the value is extracted "foodId"
		filter := bson.M{"food_id":foodId}

		// creating an instance of the options.UpdateOptions struct used
		// to specify various options
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		// updating the item in the database
		result,err := foodCollection.UpdateOne(
			// context
			ctx,
			// filter the doc to update
			filter,
			// using the set operator to set the field in "updateObj"
			bson.D{
				{Key: "$set",Value: updateObj},
			},
			// passes the options.UpdateOptions to control the update behaviour
			&opt,
		)

		// returns internal server error incase the update operation fails
		if err != nil {
			msg := fmt.Sprintf("food item update failed")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		
		// returns the result of the operation as a JSON format response
		c.JSON(http.StatusOK,result)

	}
}

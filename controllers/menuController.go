package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-backend/database"
	"restaurant-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// gin.Handler/func represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

// opens a MongoDb collection specified by the name menu
var menuCollection *mongo.Collection = database.OpenCollection(database.Client,"menu")

func GetMenus() gin.HandlerFunc{
	// Handler function for getting the menu items
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		
		// Finding all the documents in the menuCollection.
		// bson.M{} is an empty filter indicating that all documents should be retrieved
		result,err :=  menuCollection.Find(context.TODO(),bson.M{})

		// Ensures that the sorrounding context is closed when the GetMenus completes
		defer cancel()

		// Checks if an error occured during the mongodb operation. If there's the error returns a json
		// and function exits
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing the menu item"})
		}

		// Declaring a slice to store the retrieved menu items
		var allMenus []bson.M

		// Extract all documents,from the mongodb result into the  "allmenu" slice
		if err = result.All(ctx,&allMenus); err != nil {
			log.Fatal(err)
		}

		// If everything is successful, the retrieved menu items are returned as a JSON response
        c.JSON(http.StatusOK,allMenus)
	}
}

func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a time out of 100 seconds
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)

		// create a menu instance
		var menu models.Menu
		// retrieving the value of the menu id from the http request
		menuId := c.Param("menu_id")

		// Querying the database to check if there is a document with the 
		// corresponding ID
		err := foodCollection.FindOne(ctx,bson.M{"menu_id":menuId}).Decode(&menu)

		// cancel the context after the database operation
		defer cancel()

		// handle the error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while fetching the menu"})
			return
		}

		// if the operation is successful
		c.JSON(http.StatusOK,menu)

	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context) {
        // creating an instance of the menu struct
		var menu models.Menu

		// creating a context with a time out of 100 seconds
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)

		// used to extract and decode JSON data from a HTTP request body to the menu struct
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}

		// validating the data input in the menu whether it's in the right order and format
		validationErr := validate.Struct(menu)
		// Error handling incase the input data to the struct is not in the correct format
		if validationErr != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}

		// creating the time stamps of menu creation times
		menu.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		menu.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

		// creating the id for the struct to match the required id field in the document
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		// creating the menu struct into the menu collection
		result,insertErr := menuCollection.InsertOne(ctx,menu)
		if insertErr != nil{
			msg := fmt.Sprintf("menu item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		
		// cancel the context after the database insertion is done
		defer cancel()

		// return the result as a json response for the inserted meal with a status code
		// of 200
		c.JSON(http.StatusOK,result)

		// cancel the context resources and deadlines
		defer cancel()

	}
}

// creating the inTimeSpan
func inTimeSpan(start,end,check time.Time) bool {
	return start.After(time.Now()) && end.After(time.Now())
} 

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a time out of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// creating an instance of the menu struct
		var menu models.Menu

		// this method is used to extract and decode the JSON data from the HTTP request
		// body to the menu struct
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
			return
		}

		// Retrieve the value of the "menu_id" parameter from the request, c is the gin context
		menuId := c.Param("menu_id")
		
		// This line creates a MongoDB filter by constructing a map where the key
		// is "menu_id" and the value is the extracted "menuId". This filter can be used
		// in MongoDb queries to find documents where the "menu_id" field matches the extracted value
		filter := bson.M{"menu_id":menuId}

		// Declares a variable to store update operations in a Bson document
		var updateObj primitive.D

		// Checks if the start and end date in the menu struct are not nil 
		if menu.Start_date != nil && menu.End_date != nil {
			// and checks if there are in the correct timespan and cancels the operation is an error occurs
			if !inTimeSpan(*menu.Start_date,*menu.End_date,time.Now()){
				msg := "kindly retype the time"
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				// cancelling the context and exits the function
				defer cancel()
				return
			}
		}

		// Appending BSON elements to 'updateObj' for the start_date and end_date
		updateObj = append(updateObj, bson.E{Key :"start_date",Value: menu.Start_date})
		updateObj = append(updateObj, bson.E{Key :"end_date", Value: menu.End_date})

		// updating the "name" field only if 'Name' in the "menu" struct is not an empty string
		if menu.Name != ""{
			updateObj = append(updateObj, bson.E{Key :"name",Value: menu.Name})
		}

        // updating the "category" field only if 'Category' in the "menu" struct is not an empty string
		if menu.Category != ""{
			updateObj = append(updateObj, bson.E{Key :"category",Value :menu.Category})
		}

		// Updates the "Updated_at" field in the 'menu' struct with the current time
		menu.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key :"updated_at",Value: menu.Updated_at})
		
		// Sets up the 'Upsert' option to perform an upsert operation(insert a new document if a 
		// matching document is not found)
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result,err := menuCollection.UpdateOne(
			// context
			ctx,
			// specifies the document to update
			filter,
			// using the set operator to set the fields specified in 'updateObj'
			bson.D{
				{Key :"$set",Value: updateObj},
			},
			// provides the 'UpdateOptions' ('opt') for upsert behaviour
			&opt,
		)

		// Checks for errors during the update operation and returns an error message
		// if there is an error
		if err != nil{
			msg := "menu update failed"
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
		}

		// Defers the cancellation of the context until the function exits
		defer cancel()

		// Incase of a success, the function returns a Json response with the HTTP status
		// OK and the result of the operation
		c.JSON(http.StatusOK,result)
	}
}
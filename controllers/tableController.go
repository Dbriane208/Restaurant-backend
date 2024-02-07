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

// gin.HandlerFunc represent a request handler in gin
// func(ctx *gin.Context) represents the actual request handler for the routes
// [ctx *gin.Context] represents the actual parameters for the current HTTP request and response

// creating a table collection in the database
var tableCollection *mongo.Collection = database.OpenCollection(database.Client,"table")

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

	    // Retrieving the tableCollection data from the database and mapping the into the bson document
		result,err := tableCollection.Find(context.TODO(),bson.M{})
		// cancelling the resources context until the functions exits
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing table items"})
		}

		// creating a slice to store the retrieved data fromthe table collection
		var allTables []bson.M
		if err = result.All(ctx,&allTables); err != nil{
			log.Fatal(err)
		}

		// returning the table collection data in a JSON response
		c.JSON(http.StatusOK,allTables)
	}
}

func GetTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// creating an instance of the table struct
		var table models.Table
		// retrieving the table_id from the http request
		tableId := c.Param("table_id")

		// querying the database to find the document that matches the table_id
		err := tableCollection.FindOne(ctx,bson.M{"table_id":tableId}).Decode(&table)
		// cancelling the context resources until the function exits 
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error when getting a table"})
			return
		}

		// returning a JSON response for document that matched the table id
		c.JSON(http.StatusOK,table)

	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a time context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// creating an instance of the table struct
		var table models.Table

		// extracting and decoding the body of the http request into the table struct
		if err := c.BindJSON(&table); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		// validating the input data of the table to check if it's in the right format
		validationErr := validate.Struct(table)

		// Handling the error
		if validationErr != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}

		// updating the update and created at time to the current time
		table.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		table.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

		// initializing the table id in a hexadecimal representation
		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		// inserting the table document into the table collection
		result,insertErr := tableCollection.InsertOne(ctx,table)
		if insertErr != nil{
			msg := fmt.Sprintf("Table item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		// canceling the context resources until the the function exits
		defer cancel()

		// returning the created table item as a JSON response
		c.JSON(http.StatusOK,result)
	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		 
		// creating an instance of the table struct
		var table models.Table
        // retrieving the table_id from the http request
		tableId := c.Param("table_id")
        // extracting and decoding the http body to the table struct
		if err := c.BindJSON(&table); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		// creating a variable to track and store updates of the bson document
		var updateObj primitive.D

		// updating the number of guest if they're not nil
		if table.Number_of_guests != nil{
			updateObj = append(updateObj, bson.E{Key: "number_of_guests",Value: table.Number_of_guests})
		}

		// updating the table number if not nil
		if table.Table_number != nil{
			updateObj = append(updateObj, bson.E{Key: "table_number",Value: table.Table_number})
		}

		// updating the updated at time to current time
		table.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at",Value: table.Updated_at})

		// setting the operation of update and insert to true 
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		// creating a filter to select the document which should be updated
		filter := bson.M{"table_id":tableId}
        // updating the document with the corresponding ID in the database
		result,err := tableCollection.UpdateOne(
			// context
			ctx,
			// document to update
			filter,
			// updating the document using the set operator
			bson.D{
				{Key: "$set",Value: updateObj},
			},
			// providing the upsert optinks behaviour
			&opt,
		)

		// handling the error
		if err != nil{
			msg := fmt.Sprintf("table item update failed")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		// deffering the cancellation of the context resources until the function exits
		defer cancel()
		// returning the updated document result as a JSON response
		c.JSON(http.StatusOK,result)
	}
}
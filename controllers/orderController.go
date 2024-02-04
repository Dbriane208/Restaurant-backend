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

// create the order collection
var orderCollection *mongo.Collection = database.OpenCollection(database.Client,"order")


func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// querying the database to get all the items and store the in the bson map
		result,err := orderCollection.Find(context.TODO(),bson.M{})
        // canceling the resources until the function exits
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing order items"})
		}

		// creating varible of type bson slice to store the retrieved items
		var allOrders []bson.M
		// extract all documents from the mongodb result into the allOrders slice
		if err = result.All(ctx,&allOrders); err != nil{
			log.Fatal(err)
		}

		// returning a JSON response of the documents extracted from the request
		c.JSON(http.StatusOK,allOrders)

	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		 
		// retrieving the value of order_id paramenter from the http request
		orderId := c.Param("order_id")
		// creating an instance of the order struct
		var order models.Order

		// Querying the database to find the document that matches the order_id and 
		// decoding the result into the order struct
		err := foodCollection.FindOne(ctx,bson.M{"order_id":orderId}).Decode(&order)
		// cancelling the resources until the function exits
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching the orders"})
			return
		}

		// returning the JSON response with the retrieved document
		c.JSON(http.StatusOK,order)
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {

		// creating an instance of the table and the models instance
		var table models.Table
		var order models.Order

		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// creating a variable to track the updates in the bson document
		var updateObj primitive.D

		// retrieving the value of the order id from the http request
		orderId := c.Param("order_id")
		// extracting and decoding the http request body into the order struct
		if err := c.BindJSON(&order); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		if order.Table_id != nil{
			// querying the database to find the document that matches the id  table id and decoding
			// data of the result into the table struct
			err := orderCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&table)
			// cancelling the resources until the function exits
			defer cancel()
			if err != nil{
				msg := fmt.Sprintf("message: Order was not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				return
			}
			// updating the "table id" incase the  table id is not null
			updateObj = append(updateObj, bson.E{Key :"table_id",Value: order.Table_id})
		}

		// updating the "updated_at" field to the current time
		order.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at",Value: order.Updated_at})

		// setting the upsert option to true to insert the  document if it doesn't 
		// exist or update it.
		upsert := true
		// This line creates a MongoDB filter by constructing a map where the key
		// is "order_id" and the value is the extracted "orderId". This filter can be used
		// in MongoDb queries to find documents where the "menu_id" field matches the extracted value
		filter := bson.M{"order_id":orderId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		// updating the bson document using the aggregate functions
		result,err := orderCollection.UpdateOne(
			// context
			ctx,
			// specifies the document to update
			filter,
			// using the set operation to update the updateObj fields
			bson.D{
				{Key :"$set",Value: updateObj},
			},
			// provides the updateOptions for upsert behaviour
			&opt,
		)

		if err != nil{
			msg := fmt.Sprintf("order item update failed")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		// cancelling the context resources until the function exits
		defer cancel()
		// returning a JSON response with the result and status OK
		c.JSON(http.StatusOK,result)

	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating an instance of the table and order struct
		var table models.Table
		var order models.Order

		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// extracting and decoding the http request body into the order struct
		if err := c.BindJSON(&order); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		// validating the format of the input data in the order struct
		validationErr := validate.Struct(order)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}

		if order.Table_id != nil{
			// querying the document that matches the table_id and decoding the data into the table struct
			err := tableCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&table)
			// cancelling the resources and context until the function exits 
			defer cancel()
			if err != nil{
				msg := fmt.Sprintf("message:Table was not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				return
			}
		}

		// updating the time stamps of the created_at and the updated at to the current time
		order.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		order.Updated_at,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

		// Initializing the the order id and giving it hexadecimal representation
		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()

		// creating the order struct in the orderCollection
		result,insertErr := orderCollection.InsertOne(ctx,order)
		if insertErr != nil {
			msg := fmt.Sprintf("order item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		// cancelling the resources and context until the function exits
		defer cancel()

		// returns a JSON response with the inserted data to the database
		c.JSON(http.StatusOK,result)
	}
}

func orderItemOrderCreator(order models.Order) string{

	// creating a context with a timeout of 100 seconds
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

	// parsing and setting the created_at time and updated_at time field of the order
	order.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
	order.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

	// generating a unique order_id using the hexadecimal representation
	// of the order's id
	order.Order_id = order.ID.Hex()

	// inserting the order into the orderCollection in the database
	orderCollection.InsertOne(ctx,order)

	// deferring the cancellation of the context until the function returns
	defer cancel()

	// returning the generated Order_id
	return order.Order_id
}
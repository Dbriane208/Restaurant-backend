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

type orderItemsPack struct{
	Table_id *string
	Order_items []models.OrderItem
}

// creating the orderItems collection in the database
var orderItemsCollection *mongo.Collection = database.OpenCollection(database.Client,"orderItems")

func GetOrderItems() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
        // quering all the orderitem in the database
		result,err := orderItemsCollection.Find(context.TODO(),bson.M{})
		// defering the cancellation of context resources until the function exits
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing ordered items"})
			return
		}
        // creating a slice to store the retrieved items form the database
		var allOrderItems []bson.M
		// storing the items from the database to the allOrderItems
		if err = result.All(ctx,&allOrderItems); err != nil{
			log.Fatal(err)
			return
		}
		// returning the response of the retrieved items as a JSON response
		c.JSON(http.StatusOK,allOrderItems)

	}
}

func GetOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
 
		// retrieving the order_item_id from the http request
		orderItemId := c.Param("order_item_id")
		// creating an instance of the orderItem struct
		var orderItem models.OrderItem

		// querying the database to find the document that matches the order_item_id
		err := orderItemsCollection.FindOne(ctx,bson.M{"orderItem_id":orderItemId}).Decode(&orderItem)
		// cancelling the context resources until the function exits 
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing ordered item"})
			return
		}
		// returning a JSON response for document that matched the orderItem_id
		c.JSON(http.StatusOK,orderItem)

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItems,err := ItemByOrder(orderId)

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing order items by ID"})
			return
		}
		c.JSON(http.StatusOK,allOrderItems)
	}
}

func ItemByOrder(id string) (orderItems []primitive.M, err error){

} 

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var orderItem models.OrderItem

		orderItemId := c.Param("order_item_id")
		filter := bson.M{"order_item_id":orderItemId}

		var updateObj primitive.D

		orderItem.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
        updateObj = append(updateObj, bson.E{Key: "updated_at",Value: orderItem.Updated_at})

		if orderItem.Unit_price != nil{
			updateObj = append(updateObj, bson.E{Key: "unit_price",Value: &orderItem.Unit_price})
		}

		if orderItem.Quantity != nil{
			updateObj = append(updateObj, bson.E{Key: "quantity",Value: &orderItem.Quantity})
		}

		if orderItem.Food_id != nil{
			updateObj = append(updateObj, bson.E{Key: "food_id",Value: &orderItem.Food_id})
		}

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result,err := orderItemsCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set",Value: updateObj},
			},
			&opt,
		)

		if err != nil{
			msg := fmt.Sprintf("order item update failed")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK,result)
	}
}

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var orderItemPack orderItemsPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		order.Order_date,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

		orderItemToBeInserted := []interface{}{}
		order.Table_id = orderItemPack.Table_id
		order_id := orderItemOrderCreator(order)

		for _,orderItem := range orderItemPack.Order_items{
			orderItem.Order_id = order_id

			validationErr := validate.Struct(orderItem)

			if validationErr != nil{
				c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Order_item_id = orderItem.ID.Hex()

			orderItem.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
			orderItem.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

			var num = toFixed(*orderItem.Unit_price,2)
			orderItem.Unit_price = &num
			orderItemToBeInserted = append(orderItemToBeInserted, orderItem)
			
		}

		insertedOrderItem, err := orderItemsCollection.InsertMany(ctx,orderItemToBeInserted)
		if err != nil{
			log.Fatal(err)
		}
		defer cancel()

		c.JSON(http.StatusOK,insertedOrderItem)
	}
}
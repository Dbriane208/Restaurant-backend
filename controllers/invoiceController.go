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

// gin.HandlerFunc represents a request handler in gin
// func(ctx *gin.Context) represents an anonymous function which handles the actual route request
// [ctx *gin.Context] represents the current HTTP request and response

// creating a format of the invoice
type InvoiceViewFormat struct {
	Invoice_id        string
	Payment_method    string
	Order_id          string
	Payment_status   *string
	Payment_due       interface{}
	Table_number      interface{}
	Payment_due_date  time.Time
	Order_details     interface{}
}

// creating a collection od the invoices in the database
var invoicesCollection *mongo.Collection = database.OpenCollection(database.Client,"invoices")

func GetInvoices() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		 
		// Instantiating a query to the database to query all the elements and store them in a bson Map
		result,err :=  invoicesCollection.Find(context.TODO(),bson.M{})
		// cancelling all the resources until the function exits 
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while getting all the invoices"})
			return
		}
        // creating a slice variable to store the the extracted data from the database query
		var allInvoices []bson.M
		if err := result.All(ctx,&allInvoices); err != nil{
			log.Fatal(err)
		}

		// Returning the allInvoices records incase of a success
		c.JSON(http.StatusOK,allInvoices)

	}
}

func GetInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a context of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// creating a instance of the invoice
		var invoice models.Invoice
		// Retrieving the "invoice_id" from the request
		invoiceId := c.Param("invoice_id")

		// querying the data database to match the id invoice and decode the data form the request into the invoice variable
		err := invoicesCollection.FindOne(ctx,bson.M{"invoice_id":invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error when getting the invoice item"})
			return
		}

		// creating an instance of the invoiceViewformat
		var invoiceView InvoiceViewFormat

		// Retrieving order items by calling the ItemByOrder function
		// to retrieve order items based on the invoice's order ID
		allOrderItem,err := ItemByOrder(invoice.Order_id)
		if err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}
		
		// populating the fields of 'invoiceView' with data from the retrieved invoice
		// and invoice and order items
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date

		invoiceView.Payment_method = "null"
		if invoice.Payment_method != nil{
			invoice.Payment_method = *&invoice.Payment_method
		}
		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrderItem[0]["payment_due"]
		invoiceView.Table_number = allOrderItem[0]["table_number"]
		invoiceView.Order_details = allOrderItem[0]["order_item"]

		// returning JSON response with the constructed invoiceView
		c.JSON(http.StatusOK,invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a timeout of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		// creating an instance of the invoice struct
		var invoice models.Invoice

		// extracting the data from the request and decoding it in the invoice struct
		if err := c.BindJSON(&invoice); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
        // creating the instance of the order struct
		var order models.Order
		// quering the database to find the document with id -> order_id
		err := orderCollection.FindOne(ctx,bson.M{"order_id":invoice.Order_id}).Decode(&order)
		// cancelling the resources until the function returns
		defer cancel()

		if err != nil{
			msg := fmt.Sprintf("message: Order was not found")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		// declaring the status variable and setting it as the state incase the 
		// payment status is equal to nil -> success
		status := "PENDING"
		if invoice.Payment_status == nil{
			invoice.Payment_status = &status
		}

		// updating the time stamps wiht the current time
		invoice.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		invoice.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		invoice.Payment_due_date,_ = time.Parse(time.RFC3339,time.Now().AddDate(0,0,1).Format(time.RFC3339))
		
		// initalizing the id of the invoice struct and giving it the hexadecimal representation
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id =invoice.ID.Hex()

		// validating the invoice struct to check whether the data received is in the right format
		validationErr := validate.Struct(invoice)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}

		// Database operation to insert the data binded into the invoice struct
		// and handling the error incase the invoice was not created
		result,insertError := invoicesCollection.InsertOne(ctx,invoice)
		if insertError != nil{
			msg := fmt.Sprintf("invoice item was not created")
			c.JSON(http.StatusBadRequest,gin.H{"error":msg})
			return
		} 
        // cancelling all the resources until the function returns
		defer cancel()
		// returning a JSON response of the inserted data to the database
		c.JSON(http.StatusOK,result)

	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		// creating a context with a context of 100 seconds
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		// creatingg an instance of th invoice
		var invoice models.Invoice
		// Retrieving the id parameter from the request which matched id "invoice_id"
		invoiceId := c.Param("invoice_id")

		// Extracting and decoding the data from the request to the invoice struct
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		// creating a filter using the invoice id
		filter := bson.M{"invoice_id":invoiceId}
		// creatinga variable to store any updated data 
		var updateObj primitive.D

		// appending the payment_method to the updateObj variable
		if invoice.Payment_method != nil{
			updateObj = append(updateObj, bson.E{Key: "payment_method",Value: invoice.Payment_method})
		}

        // appending the payment status to the updateObj variable
		if invoice.Payment_status != nil{
			updateObj = append(updateObj, bson.E{Key: "payment_status",Value: invoice.Payment_status})
		}

		// updating the timestamp to the current time and storing the update to the updateObj variable
		invoice.Updated_at,_=time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key :"updated_at", Value: invoice.Updated_at})

		// combination of the update and insert options
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		// declaring the status variable and setting it as the state incase the 
		// payment status is equal to nil -> success
		status := "PENDING"
		if invoice.Payment_status == nil{
			invoice.Payment_status = &status
		}

		// using the aggregate function set to update the updateObj variable
		result,err := invoicesCollection.UpdateOne(
			// context
			ctx,
			// specifies the document to update
			filter,
			bson.D{
				{Key: "$set",Value: updateObj},
			},
			&opt,
		)

		if err != nil{
			msg := fmt.Sprintf("Invoice item update failed")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
 
		// deffering the cancellation of the context until the function exits
		defer cancel()
		// returns a JSON response with the result of the operation
		c.JSON(http.StatusOK,result)
		
	}
}
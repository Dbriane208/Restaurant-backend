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

var invoicesCollection *mongo.Collection = database.OpenCollection(database.Client,"invoices")

func GetInvoices() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		 
		result,err :=  invoicesCollection.Find(context.TODO(),bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while getting all the invoices"})
			return
		}

		var allInvoices []bson.M
		if err := result.All(ctx,&allInvoices); err != nil{
			log.Fatal(err)
		}

		c.JSON(http.StatusOK,allInvoices)

	}
}

func GetInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")

		err := invoicesCollection.FindOne(ctx,bson.M{"invoice_id":invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error when getting the invoice item"})
			return
		}

		var invoiceView InvoiceViewFormat

		allOrderItem,err := ItemByOrder(invoice.Order_id)
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

		c.JSON(http.StatusOK,invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		filter := bson.M{"invoice_id":invoiceId}
		var updateObj primitive.D

		if invoice.Payment_method != nil{

		}

		if invoice.Payment_status != nil{

		}

		invoice.Updated_at,_=time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key :"updated_at", Value: invoice.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		status := "PENDING"
		if invoice.Payment_status == nil{
			invoice.Payment_status = &status
		}

		result,err := invoicesCollection.UpdateOne(
			ctx,
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

		defer cancel()
		c.JSON(http.StatusOK,result)
		
	}
}
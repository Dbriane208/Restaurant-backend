package routes

// importing the necessary libraries
import (
	controller "restaurant-backend/controllers"

	"github.com/gin-gonic/gin"
)

// function used to configure routes related to Invoice operation
// function takes in an argument,incomingRoutes of type *gin.Engine
func InvoiceRoutes(incomingRoutes *gin.Engine){
	// Get request that retrieves a list of invoices
	incomingRoutes.GET("/invoices",controller.GetInvoices())
	// Get request that retrieves a specific invoice
	incomingRoutes.GET("/invoices/:invoice_id",controller.GetInvoice())
	// Post request that creates a new invoice to the database
	incomingRoutes.POST("/invoices",controller.CreateInvoice())
	// Patch request that updates a specific item entry
	incomingRoutes.PATCH("/invoices/:invoice_id",controller.UpdateInvoice())
}
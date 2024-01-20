package database

// importing necessary packages
import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// a function that returns a mongo client pointer
func DBinstance() *mongo.Client{
	// initializing a connection string and printing it out
	MongoDb := "mongodb://localhost:27017"
	fmt.Println(MongoDb)

	// creating a new mongodb client using connection string to the mongodb server 
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	// logging the error and exiting the program if an error occurs
	if err != nil {
		log.Fatal(err)
	}

	// setting a timer of 10 minutes to attempt connecting to the mongodb server
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)

	// cancelling the context if the timer is over to prevent delays
	defer cancel()

	// connecting the context to a client(mongodb server) if successful and catching an error if it fails
	err = client.Connect(ctx)
	// loggs the error and exits the program if an error occurs
	if err != nil{
		log.Fatal(err)
	}

	// returning a client if a successful connection is established
	fmt.Println("connected to mongodb")
	return client
}

// creating a global variable for the DBinstance with a mongodb client pointer
var Client *mongo.Client = DBinstance()

// function that receives two argunments a pointer and a string and returns a pointer to a collection
func OpenCollection(client *mongo.Client,collectionName string) *mongo.Collection{
	// retrieving a specified collection form the restaurant database using the provided client
    var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)
	// returning the collection
	return collection
}
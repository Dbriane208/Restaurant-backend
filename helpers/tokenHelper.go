package helpers

import (
	"context"
	"log"
	"os"
	"restaurant-backend/database"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Defines a struct with the below details and embeds jwt.
type SignedDetails struct{
	Email string
	First_name string
	Last_name string
	Uid string
	jwt.StandardClaims
}
// creates a user collection in the database
var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
// value retrieved from the environment variable
var SECRET_KEY string = os.Getenv("SECRET_KEY")

// function that takes four arguments and returns 3 values
func GenerateAllTokens(email string,firstName string,lastName string,uid string)(signedToken string,signedRefreshToken string, err error){
	// creates a variable of type *SignedDetails and initializes it with the received values
	// sets the expiry time to 24hrs from the current time
	claims := &SignedDetails{
		Email: email,
		First_name: firstName,
		Last_name: lastName,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour*time.Duration(24)).Unix(),
		},
	}

    // creates a refreshClaims variable of type *SignedDetails with only Standardclaims initalized
	// setting its expiration time to 168 hours
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour*time.Duration(168)).Unix(),
		},
	}

	// Generates a JWT token and a refresh token using the claims and refresh claims respectively signed with the 
	// HMAC-SHA256 algorithm
	token,err := jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY))

    // panics the token err generation and logs it
	if err != nil{
		log.Panic(err)
		return
	}

	// Generates a JWT token and a refresh token using the claims and refresh claims respectively signed with the 
	// HMAC-SHA256 algorithm
	refreshToken,rtErr := jwt.NewWithClaims(jwt.SigningMethodHS256,refreshClaims).SignedString([]byte(SECRET_KEY))

	// panics the token err generation and logs it
	if rtErr != nil {
		log.Panic(rtErr)
		return
	}
	
	// returns the generated token, refresh token and any error 
	return token,refreshToken,err
}

// A function that takes three arguments
func UpdateAllTokens(signedToken string,signedRefreshToken string,userid string){
	// creating a context with a timeout of 100 seconds
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

	// creating a variable to track the bson document that is changed and store the changes
	var updateObj primitive.D

	// The lines append key-value pairs to the updateobj slice. Each Key-value pairs represents an update operation
	// for the document.
	updateObj = append(updateObj, bson.E{Key: "token",Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token",Value: signedRefreshToken})

	//The line retrieves the current time formats it to RFC3339 and then parses it back to time and adds a current time
	Updated_at,_ := time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at",Value: Updated_at})

	// initializing an update and insert operation to true and using it in the update options
	upsert := true
	// creating a filter to the document we want to update which matches the id
	filter := bson.M{"user_id":userid}
	opt:= options.UpdateOptions{
		Upsert: &upsert,
	}

	// performs the update operations on the usercollection wiht all Operations in the updateObj
	// uisng the set operator
	_,err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set",Value: updateObj},
		},
		&opt,
	)
	
	// ensures that context is canceled when the function completes, releasing anu resource associeted with it
	defer cancel()

	// logs and panics the error during the update operation
	if err != nil{
		log.Panic(err)
		return
	}
}

// function that receives an argument and returns claims and a message
func ValidateToken(signedToken string)(claims *SignedDetails,msg string){
	// passing the signedToken and the signedDetails and uses an anonymous jwt token
	// to return a string slice of the secret key
	token,err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY),nil
		})	

    // asseriting the token claims to the *SignedDetails. If successful it assigns
	// the claims to the claim variable. If it fails it throws an error.
	claims,ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid :" + err.Error()
		return
	}

	// This line checks if the token has expired by comparing its expiration time with
	// the current time
	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = "Token is expired :" + err.Error()
		return
	}

	// returns the claims and the message
	return claims,msg
}
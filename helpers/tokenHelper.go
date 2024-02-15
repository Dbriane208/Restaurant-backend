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


func UpdateAllTokens(signedToken string,signedRefreshToken string,userid string){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token",Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token",Value: signedRefreshToken})
	Updated_at,_ := time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at",Value: Updated_at})

	upsert := true
	filter := bson.M{"user_id":userid}
	opt:= options.UpdateOptions{
		Upsert: &upsert,
	}

	_,err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set",Value: updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil{
		log.Panic(err)
		return
	}
}

func ValidateToken(signedToken string)(claims *SignedDetails,msg string){
	token,err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY),nil
		})	

	claims,ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid :" + err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = "Token is expired :" + err.Error()
		return
	}

	return claims,msg
}
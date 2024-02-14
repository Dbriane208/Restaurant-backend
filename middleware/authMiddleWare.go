package middleware

import (
	"net/http"
	helper "restaurant-backend/helpers"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc{
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == ""{
			msg:= "No Authorization header provide"
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			c.Abort()
			return
	     }

		 claims,err := helper.ValidateToken(clientToken)
		 if err != " "{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			c.Abort()
			return
		 }

		 c.Set("email",claims.Email)
		 c.Set("first_name",claims.First_name)
		 c.Set("last_name",claims.Last_name)
		 c.Set("uid",claims.Uid)

		 c.Next()
	}
			
}
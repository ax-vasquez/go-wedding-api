package middleware

import (
	"net/http"

	"github.com/ax-vasquez/wedding-site-api/helper"
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
)

func AuthenticateV1() gin.HandlerFunc {

	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		u := models.User{
			BaseModel: models.BaseModel{
				ID: claims.ID,
			},
		}

		// Load user from the DB to ensure context has the latest info for the user (in case there have been manual DB updates not reflected in the JWT)
		findErr := models.FindUser(c, &u)
		if findErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", u.Email)
		c.Set("first_name", u.FirstName)
		c.Set("last_name", u.LastName)
		c.Set("uid", claims.ID)
		c.Set("user_role", u.Role)
		c.Next()
	}
}

package middleware

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ax-vasquez/wedding-site-api/helper"
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
)

func AuthenticateV1() gin.HandlerFunc {

	return func(c *gin.Context) {
		test_env_str, _ := os.LookupEnv("USE_MOCK_DB")
		isUnitTest, _ := strconv.ParseBool(test_env_str)
		// If the route is being hit during unit testing, then there is no underlying DB - skip authentication since it doesn't work on unit tests
		if isUnitTest {
			// when using gin.CreateTestContextOnly, it only places values in the REQUEST context, not the gin context. Need to grab the values from the request context.
			uid := c.Request.Context().Value("uid")
			role := c.Request.Context().Value("user_role")
			c.Set("uid", uid)
			c.Set("user_role", role)
			c.Next()
			return
		}
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

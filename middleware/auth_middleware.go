package middleware

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ax-vasquez/wedding-site-api/helper"
	"github.com/gin-gonic/gin"
)

// TODO: Alter definition setup since we can't import from controllers here (invalid import cycle)
type V1_API_RESPONSE struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}

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
		clientToken := c.Request.Header.Get("auth-token")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, V1_API_RESPONSE{
				Status:  http.StatusUnauthorized,
				Message: "No authorization header provided",
			})
			c.Abort()
			return
		}

		// NOTE: If a manual change to the user has been made (for example, "GUEST" to "ADMIN") after the JWT was generated, the user should
		// sign out and back in to generate a new token with new claims (so we don't have to hit the DB for the latest data each time).
		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, V1_API_RESPONSE{
				Status:  http.StatusUnauthorized,
				Message: err,
			})
			c.Abort()
			return
		}

		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("uid", claims.ID.String())
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

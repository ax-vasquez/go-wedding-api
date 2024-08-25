package middleware

import (
	"net/http"

	"github.com/ax-vasquez/wedding-site-api/helper"
	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, V1_API_RESPONSE{
				Status:  http.StatusUnauthorized,
				Message: "not authorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

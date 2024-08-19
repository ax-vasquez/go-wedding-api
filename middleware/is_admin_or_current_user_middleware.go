package middleware

import (
	"net/http"

	"github.com/ax-vasquez/wedding-site-api/helper"
	"github.com/gin-gonic/gin"
)

func IsAdminOrCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If the current user is an admin, don't bother with UID check - admins have full access
		if err := helper.CheckUserType(c, "ADMIN"); err == nil {
			c.Next()
			return
		}
		if err := helper.MatchUserTypeToUid(c, c.GetString("uid")); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

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
		// TODO: Handle when there is no ID in the parameter (needs "or" logic to check if the ID may be in the body)
		if err := helper.MatchUserTypeToUid(c, c.Param("id")); err != nil {
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

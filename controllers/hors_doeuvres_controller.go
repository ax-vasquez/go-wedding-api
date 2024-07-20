package controllers

import "github.com/gin-gonic/gin"

// Controller to handle creating a new hors doeuvres option
//
// The route that uses this controller must be protected so that
// only site admins can use this endpoint. All other requests
// should be rejected.
func CreateHorsDoeuvresOption(c *gin.Context) {
	// TODO: Add logic to reject unauthorized requests (and certainly do not deploy until all auth logic is wired up)

}

package controllers

import (
	"net/http"
	"os"

	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/gin-gonic/gin"
)

// GetHotelRoomReservationBlockLink gets the URL for guests to use to reserve rooms from the block of rooms
func GetHotelRoomReservationBlockLink(c *gin.Context) {
	response := types.V1_API_RESPONSE_VENUE{}
	link := os.Getenv("RESERVATIONS_LINK")
	response.Status = http.StatusOK
	response.Data.Link = link
	c.JSON(http.StatusOK, response)
}

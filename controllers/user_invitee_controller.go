package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
)

// Create a user invitee
func CreateUserInvitee(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	var input models.User
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		// Invited users are considered +1s to wedding guest; they cannot invite others
		input.CanInviteOthers = false
		result, err := models.CreateUserInvitee(uint(id), input)
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Failed to insert user invitee record."
			log.Println("Error creating user invitee: ", err.Error())
		} else {
			status = http.StatusCreated
			response.Message = "Created user invitee"
			response.Data = gin.H{"records": result}
		}
	}
	response.Status = status
	c.JSON(status, response)
}

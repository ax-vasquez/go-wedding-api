package controllers

import (
	"log"
	"net/http"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Create a user invitee
func CreateUserInvitee(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	var input models.User
	idStr := c.Param("id")
	id, _ := uuid.Parse(idStr)
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		// Invited users are considered +1s to wedding guest; they cannot invite others
		input.CanInviteOthers = false
		result, err := models.CreateUserInvitee(id, input)
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Internal server error - contact server administrator."
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

// Get all invitees for the given user
func GetInviteesForUser(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	id, _ := uuid.Parse(c.Param("invitee_id"))
	status = http.StatusOK
	response.Status = status
	data, err := models.FindInviteesForUser(id)
	if err != nil {
		status = http.StatusInternalServerError
	} else {
		response.Data = gin.H{
			"invitees": data}
	}
	c.JSON(status, response)
}

// Delete an invitee for the given user
func DeleteInviteeForUser(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	invitee_id, _ := uuid.Parse(c.Param("invitee_id"))
	status = http.StatusAccepted
	response.Status = status
	result, err := models.DeleteInvitee(invitee_id)
	if err != nil {
		status = http.StatusInternalServerError
	} else {
		response.Data = gin.H{
			"records": result}
	}
	c.JSON(status, response)
}

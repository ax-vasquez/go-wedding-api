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
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	status = http.StatusOK
	response.Status = status
	response.Data = gin.H{
		"invitees": models.FindInviteesForUser(uint(id))}
	c.JSON(status, response)
}

// Delete an invitee for the given user
func DeleteInviteeForUser(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	inviter_id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	invitee_id, _ := strconv.ParseUint(c.Param("invitee_id"), 10, 64)
	status = http.StatusAccepted
	response.Status = status
	response.Data = gin.H{
		"records": models.DeleteInvitee(uint(inviter_id), uint(invitee_id))}
	c.JSON(status, response)
}

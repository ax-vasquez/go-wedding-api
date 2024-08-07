package controllers

import (
	"log"
	"net/http"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserInviteeData struct {
	Invitees []models.User `json:"users"`
}

type V1_API_RESPONSE_USER_INVITEES struct {
	V1_API_RESPONSE
	Data UserInviteeData `json:"data"`
}

// Create a user invitee
func CreateUserInvitee(c *gin.Context) {
	response := V1_API_RESPONSE_USER_INVITEES{}
	var status int
	var invitee models.User
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		status = http.StatusBadRequest
	} else {
		if err := c.ShouldBindBodyWithJSON(&invitee); err != nil {
			status = http.StatusBadRequest
			response.Message = err.Error()
		} else {
			// Invited users are considered +1s to wedding guest; they cannot invite others
			invitee.CanInviteOthers = false
			err := models.CreateUserInvitee(id, &invitee)
			if err != nil {
				status = http.StatusInternalServerError
				response.Message = "Internal server error - contact server administrator."
				log.Println("Error creating user invitee: ", err.Error())
			} else {
				status = http.StatusCreated
				response.Message = "Created user invitee"
				response.Data.Invitees = []models.User{invitee}
			}
		}
	}

	response.Status = status
	c.JSON(status, response)
}

// Get all invitees for the given user
func GetInviteesForUser(c *gin.Context) {
	response := V1_API_RESPONSE_USER_INVITEES{}
	var status int
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {

	} else {
		status = http.StatusOK
		response.Status = status
		data, err := models.FindInviteesForUser(id)
		if err != nil {
			status = http.StatusInternalServerError
		} else {
			response.Data.Invitees = data
		}
	}
	c.JSON(status, response)
}

// Delete an invitee for the given user
func DeleteInviteeForUser(c *gin.Context) {
	response := V1_API_DELETE_RESPONSE{}
	var status int
	invitee_id, _ := uuid.Parse(c.Param("invitee_id"))
	status = http.StatusAccepted
	response.Status = status
	result, err := models.DeleteInvitee(invitee_id)
	if err != nil {
		status = http.StatusInternalServerError
	} else {
		response.Data.DeletedRecords = int(*result)
	}
	c.JSON(status, response)
}

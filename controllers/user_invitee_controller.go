package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

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

// CreateUserInvitee invites a user
//
//	@Summary      invite a user
//	@Description  Invites a user for ght given user
//	@Tags         user invitee
//	@Produce      json
//	@Success      200  {object}  V1_API_RESPONSE_USER_INVITEES
//	@Failure      400  {object}  V1_API_RESPONSE_USER_INVITEES
//	@Failure      500  {object}  V1_API_RESPONSE_USER_INVITEES
//	@Param 		  user_id  path string true "Inviting user ID" Format(uuid)
//	@Router       /user/{user_id}/invite-user [post]
func CreateUserInvitee(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := V1_API_RESPONSE_USER_INVITEES{}
	var status int
	var invitee models.User
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		if err := c.ShouldBindBodyWithJSON(&invitee); err != nil {
			status = http.StatusBadRequest
			response.Message = err.Error()
		} else {
			// Invited users are considered +1s to wedding guest; they cannot invite others
			invitee.CanInviteOthers = false
			err := models.CreateUserInvitee(ctx, id, &invitee)
			if err != nil {
				status = http.StatusInternalServerError
				response.Message = "Internal server error"
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

// GetInviteesForUser gets invitees for the given user
//
//	@Summary      gets invitees for user
//	@Description  Gets invitee user data for users invited by the given inviter ID
//	@Tags         user invitee
//	@Produce      json
//	@Success      200  {object}  V1_API_RESPONSE_USER_INVITEES
//	@Failure      400  {object}  V1_API_RESPONSE_USER_INVITEES
//	@Failure      500  {object}  V1_API_RESPONSE_USER_INVITEES
//	@Param 		  user_id  path string true "Invitee search by inviting user ID" Format(uuid)
//	@Router       /user/{user_id}/invitees [get]
func GetInviteesForUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := V1_API_RESPONSE_USER_INVITEES{}
	var status int
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		status = http.StatusOK
		data, err := models.FindInviteesForUser(ctx, id)
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Internal server error"
		} else {
			response.Data.Invitees = data
		}
	}
	response.Status = status
	c.JSON(status, response)
}

// DeleteInviteeForUser deletes an invitee for the given user
//
//	@Summary      deletes an invitee for the given user
//	@Description  Deletes an invitee for the given user
//	@Tags         user invitee
//	@Produce      json
//	@Success      200  {object}  V1_API_RESPONSE_USER_INVITEES
//	@Failure      500  {object}  V1_API_RESPONSE_USER_INVITEES
//	@Param 		  inviter_id  path string true "Invitee search by inviting user ID" Format(uuid)
//	@Param 		  invitee_id  path string true "Invitee search by inviting user ID" Format(uuid)
//	@Router       /user/{inviter_id}/invitee/{invitee_id} [delete]
func DeleteInviteeForUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := V1_API_DELETE_RESPONSE{}
	var status int
	invitee_id, err := uuid.Parse(c.Param("invitee_id"))
	if err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		result, err := models.DeleteInvitee(ctx, invitee_id)
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Internal server error"
		} else {
			status = http.StatusAccepted
			response.Data.DeletedRecords = int(*result)
		}
	}
	response.Status = status
	c.JSON(status, response)
}

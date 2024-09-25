package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserInviteeInput struct {
	// The user's first name.
	FirstName string `json:"first_name" binding:"required"`
	// The user's last name.
	LastName string `json:"last_name" binding:"required"`
}

// CreateUserInvitee invites a user
//
//	@Summary      invite a user
//	@Description  Invites a user for ght given user
//	@Tags         user invitee
//	@Produce      json
//	@Success      200  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Failure      400  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Failure      500  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Param 		  user_id  path string true "Inviting user ID" Format(uuid)
//	@Router       /user/{user_id}/add-invitee [post]
func CreateUserInvitee(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_RESPONSE_USER_INVITEES{}
	var status int
	var invitee UserInviteeInput

	if err := c.ShouldBindBodyWithJSON(&invitee); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		inviterId := c.GetString("uid")
		inviterIdUUID, _ := uuid.Parse(inviterId)
		invitee := models.UserInvitee{
			InviterId: inviterIdUUID,
			FirstName: invitee.FirstName,
			LastName:  invitee.LastName,
		}
		err := models.CreateUserInvitee(&ctx, &invitee)
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Internal server error"
			log.Println("Error creating user invitee: ", err.Error())
		} else {
			status = http.StatusCreated
			response.Message = "Created user invitee"
			response.Data.Invitees = []models.UserInvitee{invitee}
		}
	}

	response.Status = status
	c.JSON(status, response)
}

// GetInviteesForLoggedInUser gets invitees for the authenticated user making the request
//
//	@Summary      gets invitees for user
//	@Description  Gets invitee user data for users invited by the given inviter ID
//	@Tags         user invitee
//	@Produce      json
//	@Success      200  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Failure      400  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Failure      500  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Param 		  user_id  path string true "Invitee search by inviting user ID" Format(uuid)
//	@Router       /user/{user_id}/invitees [get]
func GetInviteesForLoggedInUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_RESPONSE_USER_INVITEES{}
	var status int
	inviterId := c.GetString("uid")
	inviterIdUUID, _ := uuid.Parse(inviterId)
	status = http.StatusOK
	data, err := models.FindInviteesForUser(&ctx, inviterIdUUID)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
	} else {
		response.Data.Invitees = data
	}
	response.Status = status
	c.JSON(status, response)
}

// UpdateInviteeForLoggedInUser update an invitee for the logged in user
//
//	@Summary      updates an invitee for the logged in user
//	@Description  Updates an invitee for the logged in user; this will have no effect if a user attempts to update an invitee they did not add
//	@Tags         user invitee
//	@Produce      json
//	@Success      200  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Failure      500  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Param 		  id  path string true "User ID of the invitee to delete" Format(uuid)
//	@Router       /user/invitees/{id} [patch]
func UpdateInviteeForLoggedInUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_RESPONSE_USER_INVITEES{}
	var status int
	invInput := UserInviteeInput{}

	if err := c.ShouldBindBodyWithJSON(&invInput); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	inviteeId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	inviterId := c.GetString("uid")
	inviterIdUUID, _ := uuid.Parse(inviterId)
	invitee := models.UserInvitee{
		BaseModel: models.BaseModel{
			ID: inviteeId,
		},
		InviterId: inviterIdUUID,
		FirstName: invInput.FirstName,
		LastName:  invInput.LastName,
	}

	err = models.UpdateInviteeForUser(&ctx, &invitee, inviterIdUUID)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
		log.Println("Error updating user invitee: ", err.Error())
	} else {
		status = http.StatusCreated
		response.Message = "Updated user invitee"
		response.Data.Invitees = []models.UserInvitee{invitee}
	}

	response.Status = status
	c.JSON(status, response)
}

// DeleteInviteeForLoggedInUser deletes the invitee by the given ID for the logged in user.
//
//	@Summary      deletes the given invitee if the logged in user is the one who invited them
//	@Description  Deletes the given invitee if the logged in user is the one who invited them
//	@Tags         user invitee
//	@Produce      json
//	@Success      200  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Failure      500  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Param 		  id  path string true "User ID of the invitee to delete" Format(uuid)
//	@Router       /user/invitees/{id} [delete]
func DeleteInviteeForLoggedInUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_DELETE_RESPONSE{}
	var status int

	inviterId := c.GetString("uid")
	inviterIdUUID, _ := uuid.Parse(inviterId)
	inviteeId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	log.Printf("Deleting invitee... INVITEE ID: %s | USER ID: %s", inviteeId, inviterIdUUID)
	result, err := models.DeleteInviteeForUser(&ctx, inviteeId, inviterIdUUID)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
		response.Status = status
		c.JSON(status, response)
		return
	}

	status = http.StatusAccepted
	response.Data.DeletedRecords = int(*result)
	response.Status = status
	c.JSON(status, response)
}

// DeleteInvitee deletes an invitee (intended for admin use only)
//
//	@Summary      deletes an invitee, regardless the inviter
//	@Description  Deletes an invitee, regardless the inviter
//	@Tags         user invitee
//	@Produce      json
//	@Success      200  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Failure      500  {object}  types.V1_API_RESPONSE_USER_INVITEES
//	@Param 		  id  path string true "User ID of the invitee to delete" Format(uuid)
//	@Router       /invitee/{id} [delete]
func DeleteInvitee(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_DELETE_RESPONSE{}
	var status int

	inviteeId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	result, err := models.DeleteInvitee(&ctx, inviteeId)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
		response.Status = status
		c.JSON(status, response)
		return
	}

	status = http.StatusAccepted
	response.Data.DeletedRecords = int(*result)
	response.Status = status
	c.JSON(status, response)
}

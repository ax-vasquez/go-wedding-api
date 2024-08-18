package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserData struct {
	Users []models.User `json:"users"`
}

type V1_API_RESPONSE_USERS struct {
	V1_API_RESPONSE
	Data UserData `json:"data"`
}

type UpdateUserInput struct {
	ID                      uuid.UUID  `json:"id" binding:"required"`
	IsGoing                 bool       `json:"is_going"`
	FirstName               string     `json:"first_name"`
	LastName                string     `json:"last_name"`
	Email                   string     `json:"email"`
	HorsDoeuvresSelectionId *uuid.UUID `json:"hors_douevres_selection_id"`
	EntreeSelectionId       *uuid.UUID `json:"entree_selection_id"`
}

// GetUsers gets user(s) by ID(s)
//
//	@Summary      gets user(s)
//	@Description  Gets user(s) by the ID(s) in the request query string, `?ids=`
//	@Tags         user
//	@Produce      json
//	@Success      200  {object}  V1_API_RESPONSE_USERS
//	@Failure      500  {object}  V1_API_RESPONSE_USERS
//	@Param 		  ids  path string true "user search by id" Format(uuid)
//	@Router       /users [get]
func GetUsers(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := V1_API_RESPONSE_USERS{}
	var userIds []uuid.UUID
	var status int
	userIdStrings := strings.Split(c.Query("ids"), ",")
	for _, userIdStr := range userIdStrings {
		userId, _ := uuid.Parse(userIdStr)
		userIds = append(userIds, userId)
	}
	users, err := models.FindUsers(ctx, userIds)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
	} else {
		status = http.StatusOK
	}
	response.Status = status
	response.Data.Users = users
	c.JSON(status, response)
}

// CreateUser create a user
//
//	@Summary      admin-only operation to create a user
//	@Description  Creates a user with the given input and returns an array of user objects, containing the newly-created user
//	@Tags         user
//	@Accept       json
//	@Produce      json
//	@Param		  data body models.User true "The input user data (only `first_name`, `last_name` and `email` are required)"
//	@Success      201  {object}  V1_API_RESPONSE_USERS
//	@Failure      400  {object}  V1_API_RESPONSE_USERS
//	@Failure      500  {object}  V1_API_RESPONSE_USERS
//	@Router       /user [post]
func CreateUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()
	response := V1_API_RESPONSE_USERS{}
	var status int
	var input models.User
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		createUserInput := []models.User{input}
		err := models.CreateUsers(ctx, &createUserInput)
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Internal server error"
			log.Println("Error creating user: ", err.Error())
		} else {
			status = http.StatusCreated
			response.Message = "Created new user"
			response.Data.Users = createUserInput
		}
	}
	response.Status = status
	c.JSON(status, response)
}

// UpdateUser updates a user
//
//	@Summary      updates a user
//	@Description  Updates a user with the given input
//	@Tags         user
//	@Accept       json
//	@Produce      json
//	@Param		  data body models.User true "The input user update data (only `id` is required, but is not useful without setting other fields to update)"
//	@Success      202  {object}  V1_API_RESPONSE_USERS
//	@Failure      400  {object}  V1_API_RESPONSE_USERS
//	@Failure      500  {object}  V1_API_RESPONSE_USERS
//	@Router       /user [patch]
func UpdateUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := V1_API_RESPONSE_USERS{}
	var status int
	var input UpdateUserInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	u := &models.User{
		BaseModel: models.BaseModel{
			ID: input.ID,
		},
		IsGoing:                 input.IsGoing,
		FirstName:               input.FirstName,
		LastName:                input.LastName,
		Email:                   input.Email,
		HorsDoeuvresSelectionId: input.HorsDoeuvresSelectionId,
		EntreeSelectionId:       input.EntreeSelectionId,
	}
	updateErr := models.UpdateUser(ctx, u)
	setIsGoingErr := models.SetIsGoing(ctx, u)
	if updateErr != nil || setIsGoingErr != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
		response.Status = status
		c.JSON(status, response)
		if updateErr != nil {
			log.Println("Error updating user: ", updateErr.Error())
		}
		if setIsGoingErr != nil {
			log.Println("Error updating user: ", setIsGoingErr.Error())
		}
		return
	}
	status = http.StatusAccepted
	response.Message = "Updated user"
	response.Data.Users = []models.User{*u}

	response.Status = status
	c.JSON(status, response)
}

// DeleteUser delete a user
//
//	@Summary      admin-only operation to delete a user
//	@Description  Deletes an user and returns a response to indicate success or failure
//	@Tags         user
//	@Produce      json
//	@Param 		  id  path string true "User ID" Format(uuid)
//	@Success      202  {object}  V1_API_RESPONSE_USERS
//	@Failure      400  {object}  V1_API_RESPONSE_USERS
//	@Failure      500  {object}  V1_API_RESPONSE_USERS
//	@Router       /user [delete]
func DeleteUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := V1_API_DELETE_RESPONSE{}
	var status int
	id, _ := uuid.Parse(c.Param("id"))
	result, err := models.DeleteUser(ctx, id)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
		log.Println("Error deleting user: ", err.Error())
	} else {
		status = http.StatusAccepted
		response.Message = "Deleted user"
		response.Data.DeletedRecords = int(result)
	}
	response.Status = status
	c.JSON(status, response)
}

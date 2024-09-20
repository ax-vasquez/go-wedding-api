package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetLoggedInUser gets the currently logged in user
//
//	@Summary      gets the logged in user details
//	@Description  Gets the logged in user by querying for user data in the context set using JWT claims during authentication.
//	@Tags         user
//	@Produce      json
//	@Success      200  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      500  {object}  types.V1_API_RESPONSE_USERS
//	@Router       /user [get]
func GetLoggedInUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	response := types.V1_API_RESPONSE_USERS{}
	var status int

	idStr := c.GetString("uid")
	id, _ := uuid.Parse(idStr)
	users, err := models.FindUsers(ctx, []uuid.UUID{id})
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

// GetUsers gets user(s) by ID(s)
//
//	@Summary      gets user(s)
//	@Description  Gets user(s) by the ID(s) in the request query string, `?ids=`
//	@Tags         user
//	@Produce      json
//	@Success      200  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      500  {object}  types.V1_API_RESPONSE_USERS
//	@Param 		  ids  path string true "user search by id" Format(uuid)
//	@Router       /user [get]
func GetUsers(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_RESPONSE_USERS{}
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
//	@Success      201  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      400  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      500  {object}  types.V1_API_RESPONSE_USERS
//	@Router       /user [post]
func CreateUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()
	response := types.V1_API_RESPONSE_USERS{}
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

// UpdateLoggedInUser updates the logged in user
//
//	@Summary      updates the logged in user
//	@Description  Updates the logged in user with the given input
//	@Tags         user
//	@Accept       json
//	@Produce      json
//	@Param		  data body models.User true
//	@Success      202  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      400  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      500  {object}  types.V1_API_RESPONSE_USERS
//	@Router       /user [patch]
func UpdateLoggedInUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_RESPONSE_USERS{}
	var status int
	var input types.UpdateUserInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	userIdInContext, _ := c.Get("uid")
	uidStr, ok := userIdInContext.(string)
	if !ok {
		status = http.StatusInternalServerError
		response.Message = "ID in context failed type assertion (string)"
		response.Status = status
		c.JSON(status, response)
		return
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Invalid UUID detected in context."
		response.Status = status
		c.JSON(status, response)
		return
	}

	u := &models.User{
		BaseModel: models.BaseModel{
			ID: uid,
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

func AdminUpdateUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_RESPONSE_USERS{}
	var status int
	var input types.AdminUpdateUserInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = "Invalid arguments."
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
//	@Success      202  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      400  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      500  {object}  types.V1_API_RESPONSE_USERS
//	@Router       /user [delete]
func DeleteUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response := types.V1_API_DELETE_RESPONSE{}
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

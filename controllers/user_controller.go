package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ax-vasquez/wedding-site-api/helper"
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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
	IsAdmin                 bool       `json:"is_admin"`
	IsGoing                 bool       `json:"is_going"`
	CanInviteOthers         bool       `json:"can_invite_others"`
	FirstName               string     `json:"first_name"`
	LastName                string     `json:"last_name"`
	Email                   string     `json:"email"`
	HorsDoeuvresSelectionId *uuid.UUID `json:"hors_douevres_selection_id"`
	EntreeSelectionId       *uuid.UUID `json:"entree_selection_id"`
}

var validate = validator.New()

// See [jwt-in-gin-doc]
//
// [jwt-in-gin-doc]: https://www.golang.company/blog/jwt-authentication-in-golang-using-gin-web-framework
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var response V1_API_RESPONSE_USERS
		var status int

		var user models.User
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			status = http.StatusBadRequest
			response.Status = status
			response.Message = err.Error()
			c.JSON(status, response)
			return
		}

		validationErr := validate.Struct(user)
		defer cancel()
		if validationErr != nil {
			status = http.StatusBadRequest
			response.Status = status
			response.Message = validationErr.Error()
			c.JSON(status, response)
			return
		}

		count, err := models.CountUsersByEmail(ctx, &user)
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error detected while fetching user by email"})
		}

		password := helper.HashPassword(*user.PasswordHash)
		user.PasswordHash = &password

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "The mentioned E-Mail or Phone Number already exists"})
		}

		token, refreshToken, _ := helper.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.Role, user.ID.String())
		user.Token = &token
		user.RefreshToken = &refreshToken
		userSlice := []models.User{user}
		insertErr := models.CreateUsers(&userSlice)
		if insertErr != nil {
			msg := "User Details were not Saved"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		response.Data.Users = userSlice

		defer cancel()
		c.JSON(http.StatusOK, response)
	}
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
	response := V1_API_RESPONSE_USERS{}
	var userIds []uuid.UUID
	var status int
	userIdStrings := strings.Split(c.Query("ids"), ",")
	for _, userIdStr := range userIdStrings {
		userId, _ := uuid.Parse(userIdStr)
		userIds = append(userIds, userId)
	}
	users, err := models.FindUsers(userIds)
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
	response := V1_API_RESPONSE_USERS{}
	var status int
	if err := helper.CheckUserType(c, "ADMIN"); err != nil {
		status = http.StatusUnauthorized
		response.Status = status
		c.JSON(status, response)
		return
	}
	var input models.User
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		createUserInput := []models.User{input}
		err := models.CreateUsers(&createUserInput)
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
	response := V1_API_RESPONSE_USERS{}
	var status int
	var input UpdateUserInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		u := &models.User{
			BaseModel: models.BaseModel{
				ID: input.ID,
			},
			IsAdmin:                 input.IsAdmin,
			IsGoing:                 input.IsGoing,
			CanInviteOthers:         input.CanInviteOthers,
			FirstName:               input.FirstName,
			LastName:                input.LastName,
			Email:                   input.Email,
			HorsDoeuvresSelectionId: input.HorsDoeuvresSelectionId,
			EntreeSelectionId:       input.EntreeSelectionId,
		}
		updateErr := models.UpdateUser(u)
		setAdminErr := models.SetAdminPrivileges(u)
		setCanInviteErr := models.SetCanInviteOthers(u)
		setIsGoingErr := models.SetIsGoing(u)
		if updateErr != nil || setAdminErr != nil || setCanInviteErr != nil || setIsGoingErr != nil {
			status = http.StatusInternalServerError
			response.Message = "Internal server error"
			response.Status = status
			c.JSON(status, response)
			if updateErr != nil {
				log.Println("Error updating user: ", updateErr.Error())
			}
			if setAdminErr != nil {
				log.Println("Error updating user: ", setAdminErr.Error())
			}
			if setCanInviteErr != nil {
				log.Println("Error updating user: ", setCanInviteErr.Error())
			}
			if setIsGoingErr != nil {
				log.Println("Error updating user: ", setIsGoingErr.Error())
			}
			return
		}
		status = http.StatusAccepted
		response.Message = "Updated user"
		response.Data.Users = []models.User{*u}
	}
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
	response := V1_API_DELETE_RESPONSE{}
	var status int
	if err := helper.CheckUserType(c, "ADMIN"); err != nil {
		status = http.StatusUnauthorized
		response.Status = status
		c.JSON(status, response)
		return
	}
	id, _ := uuid.Parse(c.Param("id"))
	result, err := models.DeleteUser(id)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
		log.Println("Error deleting user: ", err.Error())
	} else {
		status = http.StatusAccepted
		response.Message = "Deleted user"
		response.Data.DeletedRecords = int(*result)
	}
	response.Status = status
	c.JSON(status, response)
}

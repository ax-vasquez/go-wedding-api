package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateUserInput struct {
	ID                      uuid.UUID `json:"id" binding:"required"`
	IsAdmin                 bool      `json:"is_admin"`
	IsGoing                 bool      `json:"is_going"`
	CanInviteOthers         bool      `json:"can_invite_others"`
	FirstName               string    `json:"first_name"`
	LastName                string    `json:"last_name"`
	Email                   string    `json:"email"`
	HorsDoeuvresSelectionId uuid.UUID `json:"hors_douevres_selection_id"`
	EntreeSelectionId       uuid.UUID `json:"entree_selection_id"`
}

func GetUsers(c *gin.Context) {
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
	} else {
		status = http.StatusOK
	}
	c.JSON(status, V1_API_RESPONSE{
		Status: status,
		Data: gin.H{
			"users": users}})
}

// Create a user
func CreateUsers(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	var input models.User
	createUserInput := []models.User{input}
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = "\"first_name\", \"last_name\", and \"email\" are required"
	} else {
		result, err := models.CreateUsers(&createUserInput)
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Internal server error"
			log.Println("Error creating user: ", err.Error())
		} else {
			status = http.StatusCreated
			response.Message = "Created new user"
			response.Data = gin.H{"records": result}
		}
	}
	response.Status = status
	c.JSON(status, response)
}

// Update a user
func UpdateUser(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	var input UpdateUserInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		result, err := models.UpdateUser(&models.User{
			BaseModel: models.BaseModel{
				ID: input.ID,
			},
			IsAdmin:                 input.IsAdmin,
			IsGoing:                 input.IsGoing,
			CanInviteOthers:         input.CanInviteOthers,
			FirstName:               input.FirstName,
			LastName:                input.LastName,
			Email:                   input.Email,
			HorsDoeuvresSelectionId: &input.HorsDoeuvresSelectionId,
			EntreeSelectionId:       &input.EntreeSelectionId})
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Internal server error"
			log.Println("Error updating user: ", err.Error())
		} else {
			status = http.StatusAccepted
			response.Message = "Updated user"
			response.Data = gin.H{"records": result}
		}
	}
	response.Status = status
	c.JSON(status, response)
}

// Delete a user
func DeleteUser(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	id, _ := uuid.Parse(c.Param("id"))
	result, err := models.DeleteUser(id)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
		log.Println("Error deleting user: ", err.Error())
	} else {
		status = http.StatusAccepted
		response.Message = "Deleted user"
		response.Data = gin.H{"records": result}
	}
	response.Status = status
	c.JSON(status, response)
}

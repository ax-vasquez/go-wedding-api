package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateUserInput struct {
	ID                      uint   `json:"id" binding:"required"`
	IsAdmin                 bool   `json:"is_admin"`
	IsGoing                 bool   `json:"is_going"`
	CanInviteOthers         bool   `json:"can_invite_others"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Email                   string `json:"email"`
	HorsDouevresSelectionId *uint  `json:"hors_douevres_selection_id"`
	EntreeSelectionId       *uint  `json:"entree_selection_id"`
}

func GetUsers(c *gin.Context) {
	var userIds []uint
	userIdStrings := strings.Split(c.Query("ids"), ",")
	for _, userIdStr := range userIdStrings {
		userId, _ := strconv.ParseUint(userIdStr, 10, 64)
		userIds = append(userIds, uint(userId))
	}
	users := models.FindUsers(userIds)
	c.JSON(http.StatusOK, V1_API_RESPONSE{
		Status: http.StatusOK,
		Data: gin.H{
			"users": users}})
}

// Create a user
func CreateUser(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	var input models.User
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		result, err := models.CreateUser(&input)
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Failed to insert user record."
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

func UpdateUser(c *gin.Context) {
	response := V1_API_RESPONSE{}
	var status int
	var input UpdateUserInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else {
		result, err := models.UpdateUser(&models.User{
			Model: gorm.Model{
				ID: input.ID,
			},
			IsAdmin:                 input.IsAdmin,
			IsGoing:                 input.IsGoing,
			CanInviteOthers:         input.CanInviteOthers,
			FirstName:               input.FirstName,
			LastName:                input.LastName,
			Email:                   input.Email,
			HorsDouevresSelectionId: input.HorsDouevresSelectionId,
			EntreeSelectionId:       input.EntreeSelectionId})
		if err != nil {
			status = http.StatusInternalServerError
			response.Message = "Failed to insert user record."
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

package types

import (
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type V1_API_RESPONSE struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}

type DeleteRecordResponse struct {
	DeletedRecords int `json:"deleted_records"`
}

type V1_API_DELETE_RESPONSE struct {
	V1_API_RESPONSE
	Data DeleteRecordResponse `json:"data"`
}

type AuthDetails struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type V1_API_RESPONSE_AUTH struct {
	V1_API_RESPONSE
	Data AuthDetails `json:"data"`
}

type UserData struct {
	Users []models.User `json:"users"`
}

type V1_API_RESPONSE_USERS struct {
	V1_API_RESPONSE
	Data UserData `json:"data"`
}

type UserInviteeData struct {
	Invitees []models.User `json:"users"`
}

type V1_API_RESPONSE_USER_INVITEES struct {
	V1_API_RESPONSE
	Data UserInviteeData `json:"data"`
}

type EntreeData struct {
	Entrees []models.Entree `json:"entrees"`
}

type V1_API_RESPONSE_ENTREE struct {
	V1_API_RESPONSE
	Data EntreeData `json:"data"`
}

type HorsDoeuvresData struct {
	HorsDoeuvres []models.HorsDoeuvres `json:"hors_doeuvres"`
}

type V1_API_RESPONSE_HORS_DOEUVRES struct {
	V1_API_RESPONSE
	Data HorsDoeuvresData `json:"data"`
}

type UserLoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSignupInput struct {
	UserLoginInput
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type UpdateUserInput struct {
	IsGoing                 bool       `json:"is_going"`
	FirstName               string     `json:"first_name"`
	LastName                string     `json:"last_name"`
	Email                   string     `json:"email"`
	HorsDoeuvresSelectionId *uuid.UUID `json:"hors_douevres_selection_id"`
	EntreeSelectionId       *uuid.UUID `json:"entree_selection_id"`
}

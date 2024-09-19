package controllers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ax-vasquez/wedding-site-api/helper"
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/gin-gonic/gin"
)

// Signup signs up a new user with the provided credentials
//
//	@Summary      Signs up a new user
//	@Description  Signs up a new user
//	@Tags         auth
//	@Accept       json
//	@Produce      json
//	@Param		  data body types.UserSignupInput true "Sign up details"
//	@Success      202  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      400  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      500  {object}  types.V1_API_RESPONSE_USERS
//	@Router       /signup [post]
func Signup(c *gin.Context) {
	var response types.V1_API_RESPONSE_AUTH
	var status int
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var uInput types.UserSignupInput

	if err := c.BindJSON(&uInput); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	count, err := models.CountUsersByEmail(ctx, uInput.Email)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error when checking if user exists"
		response.Status = status
		c.JSON(status, response)
		return
	}

	if count > 0 {
		status = http.StatusUnprocessableEntity
		response.Message = "A user with this email address already exists."
		response.Status = status
		c.JSON(status, response)
		return
	}

	inviteCode := os.Getenv("INVITE_CODE")
	if uInput.InviteCode != inviteCode {
		status = http.StatusUnauthorized
		response.Message = "Invalid invite code."
		response.Status = status
		c.JSON(status, response)
		return
	}

	verifyPwResult := helper.VerifyPasswordComplexity(uInput.Password, 2, 2, 2, 8)
	if !verifyPwResult.HasExpectedDigitCt || !verifyPwResult.HasExpectedSpecialCaseCt || !verifyPwResult.HasExpectedUpperCaseCt || !verifyPwResult.HasMinLength {
		var b strings.Builder
		b.WriteString("Password failed complexity requirement(s): ")
		if !verifyPwResult.HasExpectedDigitCt {
			b.WriteString("must have 2 or more digits; ")
		}
		if !verifyPwResult.HasExpectedSpecialCaseCt {
			b.WriteString("must have 2 or more special characters; ")
		}
		if !verifyPwResult.HasExpectedUpperCaseCt {
			b.WriteString("must have 2 or more capital letters; ")
		}
		if !verifyPwResult.HasMinLength {
			b.WriteString("must be at least 8 characters in length")
		}
		status = http.StatusUnprocessableEntity
		response.Status = status
		response.Message = b.String()
		c.JSON(status, response)
		return
	}

	hashedPassword := helper.HashPassword(uInput.Password)
	var newUser = models.User{}
	newUser.FirstName = uInput.FirstName
	newUser.LastName = uInput.LastName
	newUser.Email = uInput.Email
	newUser.Password = hashedPassword

	createUserInput := []models.User{newUser}
	err = models.CreateUsers(ctx, &createUserInput)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		status = http.StatusInternalServerError
		response.Message = "Internal server error while creating user"
		response.Status = status
		c.JSON(status, response)
		return
	}

	token, refreshToken, err := helper.GenerateAllTokens(createUserInput[0].Email, createUserInput[0].FirstName, createUserInput[0].LastName, createUserInput[0].Role, createUserInput[0].ID)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		status = http.StatusInternalServerError
		response.Message = "Internal server error."
		response.Status = status
		c.JSON(status, response)
		return
	}

	status = http.StatusCreated
	response.Data.Token = token
	response.Data.RefreshToken = refreshToken
	response.Status = status
	c.JSON(status, response)
}

// Login logs in a user and returns the user details for the user (if authentication is successful)
//
//	@Summary      Logs in a user
//	@Description  Logs in a user and returns the user details for the user (if authentication is successful)
//	@Tags         auth
//	@Accept       json
//	@Produce      json
//	@Param		  data body types.UserLoginInput true "Log in details"
//	@Param		  X-CSRF-Token	header	string	true "Anti CSRF token"
//	@Success      202  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      400  {object}  types.V1_API_RESPONSE_USERS
//	@Failure      500  {object}  types.V1_API_RESPONSE_USERS
//	@Router       /login [post]
func Login(c *gin.Context) {
	var response types.V1_API_RESPONSE_AUTH
	var status int
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var inputUser types.UserLoginInput
	var dbUser models.User

	if err := c.BindJSON(&inputUser); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	dbUser.Email = inputUser.Email
	// Load the user details from the DB
	err := models.FindUser(ctx, &dbUser)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		status = http.StatusInternalServerError
		response.Message = "Internal server error during user lookup"
		response.Status = status
		c.JSON(status, response)
		return
	}

	// Check password validity for user (the DB user Password is a hash, the input password is the plain-text password)
	passIsValid := helper.VerifyPassword(dbUser.Password, inputUser.Password)
	if !passIsValid {
		status = http.StatusUnauthorized
		response.Message = "Invalid credentials."
		response.Status = status
		c.JSON(status, response)
		return
	}

	// Generate new tokens for the user once we know the pass is valid
	token, refreshToken, err := helper.GenerateAllTokens(dbUser.Email, dbUser.FirstName, dbUser.LastName, dbUser.Role, dbUser.ID)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		status = http.StatusInternalServerError
		response.Message = "Internal server error."
		response.Status = status
		c.JSON(status, response)
		return
	}

	// Update signed tokens in DB for user
	err = helper.UpdateAllTokens(token, refreshToken, &dbUser)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		status = http.StatusInternalServerError
		response.Message = "Internal server error while saving auth details"
		response.Status = status
		c.JSON(status, response)
		return
	}

	status = http.StatusAccepted
	response.Status = status
	response.Data.Token = token
	response.Data.RefreshToken = refreshToken
	c.JSON(status, response)
}

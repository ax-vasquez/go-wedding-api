package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ax-vasquez/wedding-site-api/helper"
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserLoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSignupInput struct {
	UserLoginInput
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

var validate = validator.New()

// Signup signs up a new user with the provided credentials
//
//	@Summary      Signs up a new user
//	@Description  Signs up a new user
//	@Tags         auth
//	@Accept       json
//	@Produce      json
//	@Param		  data body UserSignupInput true "Sign up details"
//	@Success      201  {object}  V1_API_RESPONSE_USERS
//	@Failure      400  {object}  V1_API_RESPONSE_USERS
//	@Failure      500  {object}  V1_API_RESPONSE_USERS
//	@Router       /signup [post]
func Signup(c *gin.Context) {
	var response V1_API_RESPONSE_USERS
	var status int
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var uInput UserSignupInput

	if err := c.BindJSON(&uInput); err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	err := validate.Struct(uInput)
	if err != nil {
		status = http.StatusBadRequest
		response.Message = err.Error()
		response.Status = status
		c.JSON(status, response)
		return
	}

	count, err := models.CountUsersByEmail(ctx, uInput.Email)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Encountered an error while fetching user data for the given email."
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
		status = http.StatusInternalServerError
		response.Message = "Internal server error while creating user"
	} else {
		status = http.StatusCreated
		response.Message = "Success"
		response.Data.Users = createUserInput
	}

	response.Status = status
	c.JSON(status, response)
}

func Login(c *gin.Context) {
	var response V1_API_RESPONSE_USERS
	var status int
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var inputUser UserLoginInput
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
		status = http.StatusInternalServerError
		response.Message = "Internal server error."
		response.Status = status
		c.JSON(status, response)
		return
	}

	fmt.Println(dbUser.Password)
	fmt.Println(inputUser.Password)
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
		status = http.StatusInternalServerError
		response.Message = "Internal server error."
		response.Status = status
		c.JSON(status, response)
		return
	}

	// Update signed tokens in DB for user
	err = helper.UpdateAllTokens(token, refreshToken, &dbUser)
	if err != nil {
		status = http.StatusInternalServerError
		response.Message = "Internal server error."
	}

	status = http.StatusOK
	response.Status = status
	response.Data.Users = []models.User{dbUser}
	c.JSON(status, response)
}

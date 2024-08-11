package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Check if the user type in the context matches the role passed to it as an argument
//
// If the userType does not match the given role, an error is returned.
func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_role")
	err = nil

	if userType != role {
		err = errors.New("you are not authorised to access this resource")
	}

	return err
}

// Checks if userType is "USER" and if uid matches the userId
//
// If either contidion is unmet, an error is returned. After checking
// if the userType is "USER" and if the uid matches userId, [CheckUserType]
// is called.
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_role")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("you are not authorised to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}

// Update the user's token in the database
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {

	// TODO: Update/store JWT for user on successful authentication
	return

}

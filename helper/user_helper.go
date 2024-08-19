package helper

import (
	"errors"
	"time"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// Check if the user type in the context matches the role passed to it as an argument
//
// If the userType does not match the given role, an error is returned.
func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_role")
	err = nil

	if userType != role {
		err = errors.New("you are not authorised to access this resource")
		return err
	}

	return err
}

// Checks if userType is "USER" and if uid matches the userId
//
// If either contidion is unmet, an error is returned. After checking
// if the userType is "USER" and if the uid matches userId, [CheckUserType]
// is called.
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	uid := c.Value("uid")
	userType := c.Value("user_role")
	err = nil

	userTypeStr, ok := userType.(string)

	if !ok || uid != userId {
		err = errors.New("you are not authorised to access this resource")
		return err
	}

	err = CheckUserType(c, userTypeStr)
	return err
}

// Update the user's token in the database
func UpdateAllTokens(signedToken string, signedRefreshToken string, user *models.User) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	u := models.User{
		BaseModel: models.BaseModel{
			ID: user.ID,
		},
		Token:        signedToken,
		RefreshToken: signedRefreshToken,
	}

	err := models.UpdateUser(
		ctx,
		&u,
	)
	return err
}

package helper

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/types"
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

// MatchUserTypeToUid checks to see if the incoming request has a user_role set in the context, and if the uid matches the user ID of the
// owner for the resource being requested/modified.
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	uid := c.Value("uid")
	userType := c.Value("user_role")
	err = nil

	userTypeStr, ok := userType.(string)

	// If userId is empty at this point, it means there was no user ID present in the URL parameters
	if userId == "" {
		// Currently, the only endpoint where this logic should be possible is when a user is updating their own records (route does not have a user ID in it)
		var updateUserInput types.UpdateUserInput
		bodyBytes, bytesErr := io.ReadAll(c.Request.Body)
		err = json.Unmarshal(bodyBytes, &updateUserInput)
		if bytesErr != nil || err != nil {
			err = errors.New("bad request")
			return err
		}
		if uid != updateUserInput.ID {
			err = errors.New("you are not authorised to access this resource")
			return err
		}
	} else {
		if !ok || uid != userId {
			err = errors.New("you are not authorised to access this resource")
			return err
		}
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

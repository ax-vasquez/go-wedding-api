//go:build integration
// +build integration

package models

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_UserInvitee_Integration(t *testing.T) {
	assert := assert.New(t)
	firstUserUuid, _ := uuid.Parse(FirstUserIdStr)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	t.Run("Can find invitees for user", func(t *testing.T) {
		invitees, err := FindInviteesForUser(&ctx, firstUserUuid)
		assert.Nil(err)
		assert.NotEmpty(invitees)
		assert.Equal("Suman", invitees[0].FirstName)
	})
	t.Run("Can create user invitee", func(t *testing.T) {
		invitee := UserInvitee{
			InviterId: firstUserUuid,
			FirstName: "Billy",
			LastName:  "McTesterson",
		}
		err := CreateUserInvitee(&ctx, &invitee)
		assert.Nil(err)
		assert.NotEmpty(invitee.ID)
		assert.Equal("Billy", invitee.FirstName)
		t.Run("Can delete an invitee", func(t *testing.T) {
			result, err := DeleteInvitee(&ctx, invitee.ID)
			assert.Nil(err)
			assert.Equal(1, int(*result))
		})
	})
}

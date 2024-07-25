package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserUserInvitee(t *testing.T) {
	assert := assert.New(t)
	t.Run("Can find invitees for user", func(t *testing.T) {

	})

	t.Run("Can batch create user invitee records", func(t *testing.T) {

	})
	t.Run("Can create user invitee", func(t *testing.T) {
		invitingUserId, _ := uuid.Parse(FirstUserIdStr)
		invitee, err := CreateUserInvitee(invitingUserId, User{
			FirstName: "Billy",
			LastName:  "McTesterson",
			Email:     "a@b.com",
		})
		assert.Equal(err, nil)
		assert.NotEmpty(invitee.ID)
		t.Run("Can delete an invitee", func(t *testing.T) {
			result, err := DeleteInvitee(invitee.ID)
			assert.Equal(err, nil)
			assert.Equal(int(*result), 1)
		})
	})
}

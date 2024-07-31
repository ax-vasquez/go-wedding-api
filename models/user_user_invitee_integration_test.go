//go:build integration
// +build integration

package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserUserInvitee(t *testing.T) {
	assert := assert.New(t)
	firstUserUuid, _ := uuid.Parse(FirstUserIdStr)
	t.Run("Can find invitees for user", func(t *testing.T) {
		invitees, err := FindInviteesForUser(firstUserUuid)
		assert.Equal(nil, err)
		assert.NotEmpty(invitees)
		assert.Equal("Suman", (*invitees)[0].FirstName)
	})
	t.Run("Can create user invitee", func(t *testing.T) {
		invitee, err := CreateUserInvitee(firstUserUuid, User{
			FirstName: "Billy",
			LastName:  "McTesterson",
			Email:     "a@b.com",
		})
		assert.Equal(nil, err)
		assert.NotEmpty(invitee.ID)
		assert.Equal("Billy", invitee.FirstName)
		t.Run("Can delete an invitee", func(t *testing.T) {
			result, err := DeleteInvitee(invitee.ID)
			assert.Equal(nil, err)
			assert.Equal(1, int(*result))
		})
	})
}

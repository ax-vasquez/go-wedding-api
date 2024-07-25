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
		assert.Equal(err, nil)
		assert.NotEmpty(invitees)
		assert.Equal((*invitees)[0].FirstName, "Suman")
	})
	t.Run("Can create user invitee", func(t *testing.T) {
		invitee, err := CreateUserInvitee(firstUserUuid, User{
			FirstName: "Billy",
			LastName:  "McTesterson",
			Email:     "a@b.com",
		})
		assert.Equal(err, nil)
		assert.NotEmpty(invitee.ID)
		assert.Equal(invitee.FirstName, "Billy")
		t.Run("Can delete an invitee", func(t *testing.T) {
			result, err := DeleteInvitee(invitee.ID)
			assert.Equal(err, nil)
			assert.Equal(int(*result), 1)
		})
	})
}

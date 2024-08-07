//go:build integration
// +build integration

package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	assert := assert.New(t)
	firstUserId, _ := uuid.Parse(FirstUserIdStr)
	t.Run("Can find users", func(t *testing.T) {
		matchingUsers, err := FindUsers([]uuid.UUID{firstUserId})
		assert.Nil(err)
		assert.Equal("Rupinder", (*matchingUsers)[0].FirstName)
	})
	t.Run("Returns an empty result when no user is found", func(t *testing.T) {
		id, _ := uuid.Parse(NilUuid)
		result, err := FindUsers([]uuid.UUID{id})
		assert.Nil(err)
		assert.Empty(result)
	})
	t.Run("Can update a user", func(t *testing.T) {
		updateUser := &User{
			BaseModel: BaseModel{
				ID: firstUserId},
			FirstName: "Jimmy"}
		err := UpdateUser(updateUser)
		assert.Nil(err)
		assert.NotEmpty(updateUser.ID)
		assert.NotEqual(NilUuid, updateUser.ID)
		assert.Equal("Jimmy", updateUser.FirstName)
		assert.Equal("McNiel", updateUser.LastName)
	})
	t.Run("Can create a user", func(t *testing.T) {
		newUsers := &[]User{
			{
				FirstName: "Glizzy",
				LastName:  "Gobbler",
				Email:     "gg@gobblez.lol"}}
		err := CreateUsers(newUsers)
		assert.Nil(err)
		assert.NotEmpty((*newUsers)[0].ID)
		assert.NotEqual(NilUuid, (*newUsers)[0].ID)
		assert.Equal("Glizzy", (*newUsers)[0].FirstName)
		t.Run("Can delete a user", func(t *testing.T) {
			result, err := DeleteUser((*newUsers)[0].ID)
			assert.Nil(err)
			assert.Equal(1, int(*result))
		})
	})
}

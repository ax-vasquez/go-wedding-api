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

func Test_UserModel_Integration(t *testing.T) {
	assert := assert.New(t)
	firstUserId, _ := uuid.Parse(FirstUserIdStr)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	t.Run("Can find a set of users", func(t *testing.T) {
		matchingUsers, err := FindUsers(ctx, []uuid.UUID{firstUserId})
		assert.Nil(err)
		assert.Equal("Rupinder", matchingUsers[0].FirstName)
	})
	t.Run("Can find a single user", func(t *testing.T) {
		u := User{
			BaseModel: BaseModel{
				ID: firstUserId,
			},
		}
		err := FindUserSafe(ctx, &u)
		assert.Nil(err)
		assert.Equal("Rupinder", u.FirstName)
	})
	t.Run("Can count users by email", func(t *testing.T) {
		count, err := CountUsersByEmail(ctx, "user_1@fakedomain.com")
		assert.Nil(err)
		assert.Equal(int64(1), count)
	})
	t.Run("Returns an empty result when no user is found", func(t *testing.T) {
		id, _ := uuid.Parse(NilUuid)
		result, err := FindUsers(ctx, []uuid.UUID{id})
		assert.Nil(err)
		assert.Empty(result)
	})
	t.Run("Can update a user", func(t *testing.T) {
		updateUser := &User{
			BaseModel: BaseModel{
				ID: firstUserId},
			FirstName: "Jimmy"}
		err := UpdateUser(ctx, updateUser)
		assert.Nil(err)
		assert.NotEmpty(updateUser.ID)
		assert.NotEqual(NilUuid, updateUser.ID)
		assert.Equal("Jimmy", updateUser.FirstName)
		assert.Equal("McNiel", updateUser.LastName)
	})
	t.Run("Can create a user", func(t *testing.T) {
		newUsers := &[]User{
			{
				Role:      "GUEST",
				FirstName: "Glizzy",
				LastName:  "Gobbler",
				Email:     "gg@gobblez.lol"}}
		err := CreateUsers(ctx, newUsers)
		assert.Nil(err)
		assert.NotEmpty((*newUsers)[0].ID)
		assert.NotEqual(NilUuid, (*newUsers)[0].ID)
		assert.Equal("Glizzy", (*newUsers)[0].FirstName)
		t.Run("Can delete a user", func(t *testing.T) {
			result, err := DeleteUser(ctx, (*newUsers)[0].ID)
			assert.Nil(err)
			assert.Equal(1, int(result))
		})
	})
}

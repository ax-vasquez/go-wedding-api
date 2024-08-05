//go:build integration
// +build integration

package models

import (
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("WARNING! Could not load .env file; application will continue to run with the assumption that needed variables are present in the environment.")
	}
	os.Setenv("TEST_ENV", "true")
	Setup()
	SeedTestData()
	assert := assert.New(t)
	firstUserId, _ := uuid.Parse(FirstUserIdStr)
	t.Run("Can find users", func(t *testing.T) {
		matchingUsers, err := FindUsers([]uuid.UUID{firstUserId})
		assert.Equal(nil, err)
		assert.Equal("Rupinder", (*matchingUsers)[0].FirstName)
	})
	t.Run("Returns an empty result when no user is found", func(t *testing.T) {
		id, _ := uuid.Parse(NilUuid)
		result, err := FindUsers([]uuid.UUID{id})
		assert.Equal(nil, err)
		assert.Empty(result)
	})
	t.Run("Can update a user", func(t *testing.T) {
		updateUser, err := UpdateUser(&User{
			BaseModel: BaseModel{
				ID: firstUserId},
			FirstName: "Jimmy"})
		assert.Equal(nil, err)
		assert.Equal("Jimmy", (*updateUser)[0].FirstName)
		assert.Equal("McNiel", (*updateUser)[0].LastName)
	})
	t.Run("Can create a user", func(t *testing.T) {
		newUsers, err := CreateUsers(&[]User{
			{
				FirstName: "Glizzy",
				LastName:  "Gobbler",
				Email:     "gg@gobblez.lol"}})
		assert.Equal(nil, err)
		assert.NotEmpty((*newUsers)[0].ID)
		assert.NotEqual(NilUuid, (*newUsers)[0].ID)
		assert.Equal("Glizzy", (*newUsers)[0].FirstName)
		t.Run("Can delete a user", func(t *testing.T) {
			result, err := DeleteUser((*newUsers)[0].ID)
			assert.Equal(nil, err)
			assert.Equal(1, int(*result))
		})
	})
}

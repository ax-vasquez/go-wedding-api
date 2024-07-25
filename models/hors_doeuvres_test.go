package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHorsDoeuvres(t *testing.T) {
	assert := assert.New(t)
	t.Run("Can find a hors doeuvres", func(t *testing.T) {
		id, _ := uuid.Parse(FirstUserIdStr)
		horsDoeuvres, err := FindHorsDoeuvresById(id)
		assert.Equal(err, nil)
		assert.Equal(horsDoeuvres.OptionName, "Crab puff")
	})
	t.Run("Can find all hors doeuvres", func(t *testing.T) {
		horsDoeuvres := FindHorsDoeuvres()
		assert.Equal(len(horsDoeuvres), 5)
	})
	t.Run("Can find hors doeuvres for user", func(t *testing.T) {
		id, _ := uuid.Parse(FirstUserIdStr)
		horsDoeuvres := FindHorsDoeuvresForUser(id)
		assert.Equal(len(horsDoeuvres), 1)
		assert.Equal(horsDoeuvres[0].OptionName, "Crab puff")
	})
	t.Run("Can create an hors doeuvres", func(t *testing.T) {
		horsDoeuvresResult, err := CreateHorsDoeuvres(&[]HorsDoeuvres{{
			OptionName: "Cornflakes",
		}})
		assert.Equal(err, nil)
		assert.NotEmpty((*horsDoeuvresResult)[0].ID)
		t.Run("Can delete an hors doeuvres", func(t *testing.T) {
			id := (*horsDoeuvresResult)[0].ID
			result, err := DeleteHorsDoeuvres(id)
			assert.Equal(err, nil)
			assert.Equal(int(*result), 1)
		})
	})
}

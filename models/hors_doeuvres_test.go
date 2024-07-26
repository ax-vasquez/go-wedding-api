package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHorsDoeuvres(t *testing.T) {
	assert := assert.New(t)
	t.Run("Can find a hors doeuvres", func(t *testing.T) {
		id, _ := uuid.Parse("3baf970f-1670-4b42-ba81-63168a2f21b8")
		horsDoeuvres, err := FindHorsDoeuvresById(id)
		assert.Equal(nil, err)
		assert.Equal("Crab puff", horsDoeuvres.OptionName)
	})
	t.Run("Can find all hors doeuvres", func(t *testing.T) {
		horsDoeuvres := FindHorsDoeuvres()
		assert.Equal(5, len(horsDoeuvres))
	})
	t.Run("Can find hors doeuvres for user", func(t *testing.T) {
		id, _ := uuid.Parse(FirstUserIdStr)
		horsDoeuvres := FindHorsDoeuvresForUser(id)
		assert.Equal(1, len(horsDoeuvres))
		assert.Equal("Crab puff", horsDoeuvres[0].OptionName)
	})
	t.Run("Can create an hors doeuvres", func(t *testing.T) {
		horsDoeuvresResult, err := CreateHorsDoeuvres(&[]HorsDoeuvres{{
			OptionName: "Cornflakes",
		}})
		assert.Equal(nil, err)
		assert.NotEmpty((*horsDoeuvresResult)[0].ID)
		t.Run("Can delete an hors doeuvres", func(t *testing.T) {
			id := (*horsDoeuvresResult)[0].ID
			result, err := DeleteHorsDoeuvres(id)
			assert.Equal(nil, err)
			assert.Equal(1, int(*result))
		})
	})
}

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
		if err != nil {
			t.Fatalf("Failed fetching test entree: %s", err.Error())
		}
		assert.Equal(horsDoeuvres.OptionName, "Crab puff", "Matches expected entree option name")
	})
	t.Run("Can find all hors doeuvres", func(t *testing.T) {
		horsDoeuvres := FindHorsDoeuvres()
		assert.Equal(len(horsDoeuvres), 5, "Matches the expected number of hor doeuvres records")
	})
	t.Run("Can find hors doeuvres for user", func(t *testing.T) {
		id, _ := uuid.Parse("0ad1d80a-329b-4ffe-89c1-87af4d945953")
		horsDoeuvres := FindHorsDoeuvresForUser(id)
		assert.Equal(len(horsDoeuvres), 1, "Should only have 1 entree record in result")
		assert.Equal(horsDoeuvres[0].OptionName, "Crab puff", "Matches expected hors doeuvres option name")
	})
	t.Run("Can create an hors doeuvres", func(t *testing.T) {
		horsDoeuvresResult, err := CreateHorsDoeuvres(&[]HorsDoeuvres{{
			OptionName: "Cornflakes",
		}})
		if err != nil {
			t.Fatalf("Encountered an error while creating a new test entree: %s", err.Error())
		}
		assert.Equal((*horsDoeuvresResult)[0].OptionName, "Cornflakes", "Matches expected hors doeuvres option name")
		t.Run("Can delete an hors doeuvres", func(t *testing.T) {
			id := (*horsDoeuvresResult)[0].ID
			result, err := DeleteHorsDoeuvres(id)
			if err != nil {
				t.Fatalf("Encountered an error while deleting hors doeuvres: %s", err.Error())
			}
			assert.Equal(int(*result), 1, "Return value indicates 1 record was deleted")
		})
	})
}

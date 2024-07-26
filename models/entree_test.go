package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEntrees(t *testing.T) {
	assert := assert.New(t)
	t.Run("Can find a single entree", func(t *testing.T) {
		id, _ := uuid.Parse("f8cd5ea3-bb29-42fc-9984-a6c37d8b99c3")
		entree, err := FindEntreeById(id)
		assert.Equal(nil, err)
		assert.Equal("Caprese pasta", entree.OptionName)
	})
	t.Run("Can find all possible entrees", func(t *testing.T) {
		entrees := FindEntrees()
		assert.Equal(5, len(entrees))
	})
	t.Run("Can find entrees for user", func(t *testing.T) {
		id, _ := uuid.Parse(FirstUserIdStr)
		entrees := FindEntreesForUser(id)
		assert.Equal(1, len(entrees))
		assert.Equal("Caprese pasta", entrees[0].OptionName)
	})
	t.Run("Can create an entree", func(t *testing.T) {
		entreesResult, err := CreateEntrees(&[]Entree{{
			OptionName: "Cap'n Crunch",
		}})
		assert.Equal(nil, err)
		assert.Equal("Cap'n Crunch", (*entreesResult)[0].OptionName)
		// Embedded test so we can easily-target the new record and delete it as part of the next test
		t.Run("Can delete an entree", func(t *testing.T) {
			id := (*entreesResult)[0].ID
			result, err := DeleteEntree(id)
			assert.Equal(nil, err)
			assert.Equal(1, int(*result))
		})
	})
}

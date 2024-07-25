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
		if err != nil {
			t.Fatalf("Failed fetching test entree: %s", err.Error())
		}
		assert.Equal(entree.OptionName, "Caprese pasta")
	})
	t.Run("Can find all possible entrees", func(t *testing.T) {
		entrees := FindEntrees()
		assert.Equal(len(entrees), 5)
	})
	t.Run("Can find entrees for user", func(t *testing.T) {
		id, _ := uuid.Parse(FirstUserIdStr)
		entrees := FindEntreesForUser(id)
		assert.Equal(len(entrees), 1)
		assert.Equal(entrees[0].OptionName, "Caprese pasta")
	})
	t.Run("Can create an entree", func(t *testing.T) {
		entreesResult, err := CreateEntrees(&[]Entree{{
			OptionName: "Cap'n Crunch",
		}})
		if err != nil {
			t.Fatalf("Encountered an error while creating a new test entree: %s", err.Error())
		}
		assert.Equal((*entreesResult)[0].OptionName, "Cap'n Crunch")
		// Embedded test so we can easily-target the new record and delete it as part of the next test
		t.Run("Can delete an entree", func(t *testing.T) {
			id := (*entreesResult)[0].ID
			result, err := DeleteEntree(id)
			if err != nil {
				t.Fatalf("Encountered an error while deleting entree: %s", err.Error())
			}
			assert.Equal(int(*result), 1)
		})
	})
}

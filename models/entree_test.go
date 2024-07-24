package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestEntrees(t *testing.T) {
	t.Run("Can find a single entree", func(t *testing.T) {
		id, _ := uuid.Parse("f8cd5ea3-bb29-42fc-9984-a6c37d8b99c3")
		entree, err := FindEntreeById(id)
		if err != nil {
			t.Fatalf("Failed fetching test entree: %s", err.Error())
		}
		if entree.OptionName != "Caprese pasta" {
			t.Fatalf("Unexpected option selection found - expected \"Caprese pasta\", received: %s", entree.OptionName)
		}
	})
	t.Run("Can find all possible entrees", func(t *testing.T) {
		entrees := FindEntrees()
		if len(entrees) != 5 {
			t.Fatalf("Expected 5 entree records, received: %d", len(entrees))
		}
	})
	t.Run("Can find entrees for user", func(t *testing.T) {
		id, _ := uuid.Parse("0ad1d80a-329b-4ffe-89c1-87af4d945953")
		entrees := FindEntreesForUser(id)
		if len(entrees) > 1 {
			t.Fatalf("Users should only have 1 entree record; received: %d", len(entrees))
		}
		// All test users and invitees are set to the same initial entree selection
		if entrees[0].OptionName != "Caprese pasta" {
			t.Fatalf("Unexpected option selection found - expected \"Caprese pasta\", received: %s", entrees[0].OptionName)
		}
	})
	t.Run("Can create entrees", func(t *testing.T) {
		entreesResult, err := CreateEntrees(&[]Entree{{
			OptionName: "Cap'n Crunch",
		}})
		if err != nil {
			t.Fatalf("Encountered an error while creating a new test entree: %s", err.Error())
		}
		if (*entreesResult)[0].OptionName != "Cap'n Crunch" {
			t.Fatalf("Unexpected option name; expected \"Cap'n Crunch\", received: %s", (*entreesResult)[0].OptionName)
		}
	})
	t.Run("Can delete an entree", func(t *testing.T) {

	})
}

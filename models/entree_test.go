package models

import (
	"fmt"
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
	// TODO: Fix bug with setup/teardown - INSERT is not working properly for hors doeuvres, which puts the test DB in a bad state.
	// During the first run (while the bad state is created), this test will fail. After that, the test will pass 100% of the time.
	t.Run("Can find entrees for user", func(t *testing.T) {
		entrees := FindEntreesForUser(1)
		if len(entrees) > 1 {
			t.Fatalf("Users should only have 1 entree record; received: %d", len(entrees))
		}
		// All test users and invitees are set to the same initial entree selection
		if entrees[0].OptionName != "Caprese pasta" {
			t.Fatalf("Unexpected option selection found - expected \"Caprese pasta\", received: %s", entrees[0].OptionName)
		}
	})
	t.Run("Can create entrees", func(t *testing.T) {
		result, err := CreateEntrees(&[]Entree{{
			OptionName: "Cap'n Crunch",
		}})
		if err != nil {
			t.Fatalf("Encountered an error while creating a new test entree: %s", err.Error())
		}
		fmt.Println("RESULT: ", result)
	})
	t.Run("Can delete an entree", func(t *testing.T) {

	})
}

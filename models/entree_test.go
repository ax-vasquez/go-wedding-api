package models

import (
	"testing"
)

func TestEntrees(t *testing.T) {
	t.Run("Can find all possible entrees", func(t *testing.T) {
		entrees := FindEntrees()
		if len(entrees) != 5 {
			t.Fatalf("Expected 5 entree records, received: %d", len(entrees))
		}
	})
	// TODO: Investigate possible flake; this test failed with an index error 1-2 times, but passes all other times
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
}

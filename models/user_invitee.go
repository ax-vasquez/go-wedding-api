package models

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type UserInvitee struct {
	BaseModel
	InviterId uuid.UUID `gorm:"index" json:"inviter_id" binding:"required"`
	Inviter   User      `gorm:"foreignKey:InviterId"`
	// The user's first name.
	FirstName string `json:"first_name" binding:"required"`
	// The user's last name.
	LastName string `json:"last_name" binding:"required"`
	// The ID of the hors doeuvres the user has selected; is null until the user makes a selection.
	HorsDoeuvresSelectionId *uuid.UUID    `json:"hors_doeuvres_selection_id"`
	HorsDoeuvresSelection   *HorsDoeuvres `gorm:"foreignKey:HorsDoeuvresSelectionId"`
	// The ID of the entree the user has selected; is null until the user makes a selection.
	EntreeSelectionId *uuid.UUID `json:"entree_selection_id"`
	EntreeSelection   *Entree    `gorm:"foreignKey:EntreeSelectionId"`
}

// Create user Invitee and return the number of rows affected
//
// This inserts a new row in the user_user_invitees table, which facilitates a many-to-many relationship
// between invitee.
func CreateUserInvitee(c *context.Context, invitedUser *UserInvitee) error {
	result := db.WithContext(*c).Create(invitedUser)
	if result.Error != nil {
		log.Println("Error creating UserInvitee record: ", result.Error.Error())
		return result.Error
	}
	return nil
}

func CreateUserInvitees(c context.Context, invitees *[]UserInvitee) error {
	result := db.WithContext(c).Create(&invitees)
	return result.Error
}

// Delete an invitee
//
// This will delete the related records from the user_user_invitees table as well as the invited user from the
// users table.
func DeleteInvitee(c *context.Context, inviteeId uuid.UUID) (*int64, error) {
	result := db.WithContext(*c).Delete(&UserInvitee{}, "id = ?", inviteeId)
	if result.Error != nil {
		log.Println("Error deleting UserInvitee: ", result.Error.Error())
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

// Finds all users for the given inviting user ID
//
// It's safe to return all fields for invitees since they don't register themselves and have no auth information.
func FindInviteesForUser(c *context.Context, inviterId uuid.UUID) ([]UserInvitee, error) {
	var users []UserInvitee
	result := db.WithContext(*c).Find(&users, "inviter_id = ?", inviterId)
	if result.Error != nil {
		log.Println("Error querying for UserUserInvitee: ", result.Error.Error())
		return nil, result.Error
	}
	return users, nil
}

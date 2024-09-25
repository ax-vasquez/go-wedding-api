package models

import (
	"log"

	"github.com/google/uuid"
)

// User-UserInvitee relation table
type UserUserInvitee struct {
	BaseModel
	InviterId uuid.UUID    `gorm:"index" json:"inviter_id" binding:"required"`
	Inviter   User         `gorm:"foreignKey:InviterId"`
	InviteeId uuid.UUID    `gorm:"index" json:"invitee_id" binding:"required"`
	Invitee   *UserInvitee `gorm:"foreignKey:InviteeId"`
}

// Helper to bulk-insert multiple invitee records for users (inviter and invitee) that already exist in the database.
//
// This method is currently only used for testing purposes, and as a result, is intentionally not covered
// by any testing logic. Since this method requires the corresponding invitee User record to already exist,
// it's not useful in normal application usage. It's only useful when moving large sets of pre-existing data
// around such as when hydrating a test database with fake data. It's possible this may be useful later if
// we move data between providers, but that use is far removed from our current needs.
func CreateUserUserInvitees(u []UserUserInvitee) error {
	result := db.Create(u)
	if result.Error != nil {
		log.Println("Error creating UserUserInvitee record: ", result.Error.Error())
		return result.Error
	}
	return nil
}

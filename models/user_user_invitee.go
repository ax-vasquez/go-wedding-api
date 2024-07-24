package models

import (
	"log"

	"github.com/google/uuid"
)

// User-UserInvitee relation table
type UserUserInvitee struct {
	BaseModel
	InviterId uuid.UUID `gorm:"index" json:"inviter_id" binding:"required"`
	Inviter   User      `gorm:"foreignKey:InviterId"`
	InviteeId uuid.UUID `gorm:"index" json:"invitee_id" binding:"required"`
	Invitee   User      `gorm:"foreignKey:InviteeId"`
}

// Create user Invitee and return the number of rows affected
//
// This inserts a new row in the user_user_invitees table, which facilitates a many-to-many relationship
// between invitee
func CreateUserInvitee(inviting_user_id uuid.UUID, invited_user User) (*int64, error) {
	result := db.Create(&UserUserInvitee{
		InviterId: inviting_user_id,
		Invitee:   invited_user,
	})
	if result.Error != nil {
		log.Println("Error creating UserUserInvitee record: ", result.Error.Error())
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

// Finds all users for the given inviting user ID
func FindInviteesForUser(user_id uint) []User {
	var users []User
	result := db.Joins("JOIN user_user_invitees ON user_user_invitees.invitee_id = users.id AND user_user_invitees.inviter_id = ?", user_id).Find(&users)
	if result.Error != nil {
		log.Println("Error querying for UserUserInvitee: ", result.Error.Error())
	}
	return users
}

// Delete an invitee
//
// This will delete the related records from the user_user_invitees table as well as the invited user from the
// users table.
func DeleteInvitee(invitee_id uint) int64 {
	result := db.Delete(&UserUserInvitee{}, "invitee_id = ?", invitee_id)
	if result.Error != nil {
		log.Println("Error deleting UserUserInvitee: ", result.Error.Error())
	}
	result = db.Delete(&User{}, invitee_id)
	if result.Error != nil {
		log.Println("Error deleting invited User: ", result.Error.Error())
	}
	return result.RowsAffected
}

// Helper to bulk-insert multiple invitee records. Intended for testing, but could
// have uses elsewhere.
func CreateUserUserInvitees(u []UserUserInvitee) error {
	result := db.Create(u)
	if result.Error != nil {
		log.Println("Error creating UserUserInvitee record: ", result.Error.Error())
		return result.Error
	}
	return nil
}

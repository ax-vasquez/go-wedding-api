package models

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// User-UserInvitee relation table
type UserUserInvitee struct {
	BaseModel
	InviterId uuid.UUID `gorm:"index" json:"inviter_id" binding:"required"`
	Inviter   User      `gorm:"foreignKey:InviterId"`
	InviteeId uuid.UUID `gorm:"index" json:"invitee_id" binding:"required"`
	Invitee   *User     `gorm:"foreignKey:InviteeId"`
}

// Create user Invitee and return the number of rows affected
//
// This inserts a new row in the user_user_invitees table, which facilitates a many-to-many relationship
// between invitee.
func CreateUserInvitee(c context.Context, invitingUserId uuid.UUID, invitedUser *User) error {
	result := db.Create(&UserUserInvitee{
		InviterId: invitingUserId,
		Invitee:   invitedUser,
	})
	if result.Error != nil {
		log.Println("Error creating UserUserInvitee record: ", result.Error.Error())
		return result.Error
	}
	return nil
}

// Finds all users for the given inviting user ID
func FindInviteesForUser(c context.Context, userId uuid.UUID) ([]User, error) {
	var users []User
	result := db.WithContext(c).Joins("JOIN user_user_invitees ON user_user_invitees.invitee_id = users.id AND user_user_invitees.inviter_id = ?", userId).Find(&users)
	if result.Error != nil {
		log.Println("Error querying for UserUserInvitee: ", result.Error.Error())
		return nil, result.Error
	}
	return users, nil
}

// Delete an invitee
//
// This will delete the related records from the user_user_invitees table as well as the invited user from the
// users table.
func DeleteInvitee(c context.Context, inviteeId uuid.UUID) (*int64, error) {
	result := db.WithContext(c).Delete(&UserUserInvitee{}, "invitee_id = ?", inviteeId)
	if result.Error != nil {
		log.Println("Error deleting UserUserInvitee: ", result.Error.Error())
		return nil, result.Error
	}
	result = db.Delete(&User{}, inviteeId)
	if result.Error != nil {
		log.Println("Error deleting invited User: ", result.Error.Error())
		return nil, result.Error
	}
	return &result.RowsAffected, nil
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

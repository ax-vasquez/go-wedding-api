package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// User-UserInvitee relation table
type UserUserInvitee struct {
	gorm.Model
	CreatedAt time.Time `gorm:"<-:create"`
	InviterId uint      `gorm:"index" json:"inviter_id" binding:"required"`
	Inviter   User      `gorm:"foreignKey:InviterId"`
	InviteeId uint      `gorm:"index" json:"invitee_id" binding:"required"`
	Invitee   User      `gorm:"foreignKey:InviteeId"`
}

// Create user Invitee and return the number of rows affected
//
// This inserts a new row in the user_user_invitees table, which facilitates a many-to-many relationship
// between invitee
func CreateUserInvitee(inviting_user_id uint, invited_user User) (*int64, error) {
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

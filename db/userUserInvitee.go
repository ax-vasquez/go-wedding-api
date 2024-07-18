package db

import "gorm.io/gorm"

// User-UserInvitee relation table
//
// May not be necessary, but leads to the creation of a many-to-many table
// to link users to their invitees
type UserUserInvitee struct {
	gorm.Model
	InviterId uint
	Inviter   User `gorm:"foreignKey:InviterId"`
	InviteeId uint
	Invitee   User `gorm:"foreignKey:InviteeId"`
}

package db

import "gorm.io/gorm"

// UserInvitee table
type UserInvitee struct {
	gorm.Model
	User
	InviterId uint
	Inviter   User `gorm:"foreignKey:InviterId"`
}

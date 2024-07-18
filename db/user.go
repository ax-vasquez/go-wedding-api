package db

import "gorm.io/gorm"

// User table
type User struct {
	gorm.Model
	IsAdmin                 bool   `json:"is_admin"`
	IsGoing                 bool   `json:"is_going"`
	CanInviteOthers         bool   `json:"can_invite_others"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Email                   string `json:"email"`
	HorsDouevresSelectionId uint
	HorsDouevresSelection   HorsDouevres `gorm:"foreignKey:HorsDouevresSelectionId"`
	EntreeSelectionId       uint
	EntreeSelection         Entree `gorm:"foreignKey:EntreeSelectionId"`
}

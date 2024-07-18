package models

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// User table
type User struct {
	gorm.Model
	IsAdmin                 bool   `json:"is_admin"`
	IsGoing                 bool   `json:"is_going"`
	CanInviteOthers         bool   `json:"can_invite_others"`
	FirstName               string `json:"first_name" binding:"required"`
	LastName                string `json:"last_name" binding:"required"`
	Email                   string `json:"email" binding:"required"`
	HorsDouevresSelectionId *uint
	HorsDouevresSelection   *HorsDouevres `gorm:"foreignKey:HorsDouevresSelectionId"`
	EntreeSelectionId       *uint
	EntreeSelection         *Entree `gorm:"foreignKey:EntreeSelectionId"`
}

// Maybe create a user (if no errors) and returns the number of inserted records
func CreateUser(u *User) (int64, error) {
	result := db.Create(&u)
	if result.Error != nil {
		// Return 0 as the ID when no insert was performed
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func FindUsers(ids []uint) []User {
	var users []User
	fmt.Println(ids)
	result := db.Find(&users, ids)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
	}
	return users
}

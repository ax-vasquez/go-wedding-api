package models

import (
	"log"
	"time"

	"github.com/google/uuid"
)

// User table
type User struct {
	BaseModel
	// We override Gorm's CreatedAt field so we can set the gorm:"<-:create" directive,
	// which prevents this field from being altered once the record is created
	CreatedAt               time.Time     `gorm:"<-:create"`
	IsAdmin                 bool          `json:"is_admin"`
	IsGoing                 bool          `json:"is_going"`
	CanInviteOthers         bool          `json:"can_invite_others"`
	FirstName               string        `json:"first_name" binding:"required"`
	LastName                string        `json:"last_name" binding:"required"`
	Email                   string        `json:"email" gorm:"uniqueIndex" binding:"required"`
	HorsDoeuvresSelectionId *uuid.UUID    `json:"hors_doeuvres_selection_id"`
	HorsDoeuvresSelection   *HorsDoeuvres `gorm:"foreignKey:HorsDoeuvresSelectionId"`
	EntreeSelectionId       *uuid.UUID    `json:"entree_selection_id"`
	EntreeSelection         *Entree       `gorm:"foreignKey:EntreeSelectionId"`
}

// Maybe create users with given data (if no errors) and returns the number of inserted records
func CreateUsers(users *[]User) (*int64, error) {
	result := db.Create(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

// Maybe update a user (if no errors) and returns the number of inserted records
func UpdateUser(u *User) (*int64, error) {
	result := db.Updates(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

// Maybe delete a user (if no errors) and returns the number of deleted records
func DeleteUser(id uint) (*int64, error) {
	// Since our models have DeletedAt set, this makes Gorm "soft delete" records on normal delete operations.
	// We can add .Unscoped() prior to the .Delete() call if we want to permanently-delete them.
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

// Find Users by the given ids; returns a User slice
func FindUsers(ids []uint) []User {
	var users []User
	result := db.Find(&users, ids)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
	}
	return users
}

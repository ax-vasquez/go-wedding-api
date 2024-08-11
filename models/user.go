package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// User table
type User struct {
	BaseModel
	// We override Gorm's CreatedAt field so we can set the gorm:"<-:create" directive,
	// which prevents this field from being altered once the record is created
	CreatedAt time.Time `gorm:"<-:create"`
	// GUEST, INVITEE or ADMIN
	Role                    string `json:"role" gorm:"default:GUEST"`
	IsAdmin                 bool   `json:"is_admin"`
	IsGoing                 bool   `json:"is_going"`
	CanInviteOthers         bool   `json:"can_invite_others"`
	FirstName               string `json:"first_name" binding:"required"`
	LastName                string `json:"last_name" binding:"required"`
	Email                   string `json:"email" gorm:"uniqueIndex" binding:"required"`
	PasswordHash            *string
	Token                   *string       `json:"token"`
	RefreshToken            *string       `json:"refresh_token"`
	HorsDoeuvresSelectionId *uuid.UUID    `json:"hors_doeuvres_selection_id"`
	HorsDoeuvresSelection   *HorsDoeuvres `gorm:"foreignKey:HorsDoeuvresSelectionId"`
	EntreeSelectionId       *uuid.UUID    `json:"entree_selection_id"`
	EntreeSelection         *Entree       `gorm:"foreignKey:EntreeSelectionId"`
}

// Maybe create users with given data (if no errors) and returns the number of inserted records
func CreateUsers(users *[]User) error {
	result := db.Create(&users)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Get the count of users whose email matches that of the given user.
//
// This should only return 1 or 0 and is used to check if a user already
// exists with the given email address.
func CountUsersByEmail(user *User) (int64, error) {
	var count int64
	result := db.Distinct("email").Count(&count).Find(&user)
	return result.RowsAffected, result.Error
}

// Set is_admin for user
func SetAdminPrivileges(u *User) error {
	result := db.Model(&u).Select("is_admin").Updates(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Set is_going for user
func SetIsGoing(u *User) error {
	result := db.Model(&u).Select("is_going").Updates(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Set can_invite_others for user
func SetCanInviteOthers(u *User) error {
	result := db.Model(&u).Select("can_invite_others").Updates(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Maybe update a user (if no errors) and returns the number of inserted records
//
// The Updates() method will only update non-zero fields. Importantly, this means that
// you cannot use Updates() to set a boolean field to `false`, unless you either pass
// the updated fields as a string map, or select the fields you intend to target. However,
// this will lead to values being overwritten, thus invalidating the purpose of using
// Updates() in the first place. Use helper methods such as SetAdminPrivileges to set
// a given boolean field without overwriting unspecified fields.
//
// See: https://gorm.io/docs/update.html#Updates-multiple-columns
func UpdateUser(u *User) error {
	result := db.Model(&u).Clauses(clause.Returning{}).Updates(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Maybe delete a user (if no errors) and returns the number of deleted records
func DeleteUser(id uuid.UUID) (*int64, error) {
	// Since our models have DeletedAt set, this makes Gorm "soft delete" records on normal delete operations.
	// We can add .Unscoped() prior to the .Delete() call if we want to permanently-delete them.
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

// Find Users by the given ids; returns a User slice
func FindUsers(ids []uuid.UUID) ([]User, error) {
	var users []User
	result := db.Find(&users, ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

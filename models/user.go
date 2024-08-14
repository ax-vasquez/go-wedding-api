package models

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// User table
type User struct {
	BaseModel
	// The user's role, which can be "GUEST", "INVITEE" or "ADMIN". Defaults to "GUEST".
	Role string `json:"role" gorm:"default:GUEST"`
	// Whether or not the user is attending.
	IsGoing bool `json:"is_going"`
	// The user's first name.
	FirstName string `json:"first_name" binding:"required"`
	// The user's last name.
	LastName string `json:"last_name" binding:"required"`
	// The user's email (must be unique); this field is an index.
	Email string `json:"email" gorm:"uniqueIndex" binding:"required"`
	// The hash of the user's password.
	PasswordHash string
	// The user's auth token.
	Token string `json:"token"`
	// The user's auth refresh token.
	RefreshToken string `json:"refresh_token"`
	// The ID of the hors doeuvres the user has selected; is null until the user makes a selection.
	HorsDoeuvresSelectionId *uuid.UUID    `json:"hors_doeuvres_selection_id"`
	HorsDoeuvresSelection   *HorsDoeuvres `gorm:"foreignKey:HorsDoeuvresSelectionId"`
	// The ID of the entree the user has selected; is null until the user makes a selection.
	EntreeSelectionId *uuid.UUID `json:"entree_selection_id"`
	EntreeSelection   *Entree    `gorm:"foreignKey:EntreeSelectionId"`
}

// Maybe create users with given data (if no errors) and returns the number of inserted records
func CreateUsers(c context.Context, users *[]User) error {
	result := db.WithContext(c).Create(&users)
	return result.Error
}

// Get the count of users whose email matches that of the given user.
//
// This should only return 1 or 0 and is used to check if a user already
// exists with the given email address.
func CountUsersByEmail(c context.Context, user *User) (int64, error) {
	var count int64
	result := db.Model(&user).WithContext(c).Distinct("email").Count(&count)
	return count, result.Error
}

// Set is_going for user
func SetIsGoing(c context.Context, u *User) error {
	result := db.WithContext(c).Model(&u).Select("is_going").Updates(&u)
	return result.Error
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
func UpdateUser(c context.Context, u *User) error {
	result := db.WithContext(c).Model(&u).Clauses(clause.Returning{}).Updates(&u)
	return result.Error
}

// Maybe delete a user (if no errors) and returns the number of deleted records
func DeleteUser(c context.Context, id uuid.UUID) (*int64, error) {
	// Since our models have DeletedAt set, this makes Gorm "soft delete" records on normal delete operations.
	// We can add .Unscoped() prior to the .Delete() call if we want to permanently-delete them.
	result := db.WithContext(c).Delete(&User{}, id)
	return &result.RowsAffected, result.Error
}

// Find Users by the given ids; returns a User slice
func FindUsers(c context.Context, ids []uuid.UUID) ([]User, error) {
	var users []User
	result := db.WithContext(c).Find(&users, ids)
	return users, result.Error
}

// Find user with the given details
func FindUser(c context.Context, u *User) error {
	result := db.WithContext(c).Find(&u)
	return result.Error
}

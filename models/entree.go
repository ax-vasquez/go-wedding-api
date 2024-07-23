package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// Entree table
type Entree struct {
	gorm.Model
	CreatedAt  time.Time `gorm:"<-:create"`
	OptionName string    `json:"option_name" binding:"required"`
}

// Finds all hors doeuvres
func FindEntrees() []Entree {
	var entrees []Entree
	result := db.Find(&entrees)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
	}
	return entrees
}

// Finds hors doeuvres for the given user
func FindEntreesForUser(id uint) []Entree {
	var entrees []Entree
	result := db.Joins("JOIN users ON entrees.id = users.entree_selection_id AND users.id = ?", id).Find(&entrees)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
	}
	return entrees
}

// Maybe create a user (if no errors) and returns the number of inserted records
func CreateEntrees(entrees *[]Entree) (*int64, error) {
	result := db.Create(&entrees)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

// Maybe delete a user (if no errors) and returns the number of deleted records
func DeleteEntree(id uint) (*int64, error) {
	// Since our models have DeletedAt set, this makes Gorm "soft delete" records on normal delete operations.
	// We can add .Unscoped() prior to the .Delete() call if we want to permanently-delete them.
	result := db.Delete(&Entree{}, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

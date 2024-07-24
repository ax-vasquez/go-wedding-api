package models

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// Entree table
type Entree struct {
	BaseModel
	OptionName string `json:"option_name" binding:"required"`
}

// Finds all entrees
func FindEntrees() []Entree {
	var entrees []Entree
	result := db.Find(&entrees)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
	}
	return entrees
}

// Find a single entree by ID
func FindEntreeById(id uuid.UUID) (*Entree, error) {
	var entree *Entree
	result := db.Find(&entree, Entree{BaseModel: BaseModel{ID: id}})
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
		return nil, result.Error
	}
	return entree, nil
}

// Finds entrees for the given user
func FindEntreesForUser(id uint) []Entree {
	var entrees []Entree
	result := db.Joins("JOIN users ON entrees.id = users.entree_selection_id AND users.id = ?", id).Find(&entrees)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
	}
	return entrees
}

// Maybe create a user (if no errors) and returns the number of inserted records
func CreateEntrees(entrees *[]Entree) (*[]Entree, error) {
	result := db.Clauses(clause.Returning{}).Select("*").Create(&entrees)
	if result.Error != nil {
		return nil, result.Error
	}
	return entrees, nil
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

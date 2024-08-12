package models

import (
	"context"
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
func FindEntrees(c context.Context) ([]Entree, error) {
	var entrees []Entree
	result := db.WithContext(c).Find(&entrees)
	return entrees, result.Error
}

// Find a single entree by ID
func FindEntreeById(c context.Context, id uuid.UUID) (*Entree, error) {
	var entree *Entree
	result := db.WithContext(c).Find(&entree, Entree{BaseModel: BaseModel{ID: id}})
	return entree, result.Error
}

// Finds entrees for the given user
func FindEntreesForUser(c context.Context, id uuid.UUID) ([]Entree, error) {
	var entrees []Entree
	result := db.WithContext(c).Joins("JOIN users ON entrees.id = users.entree_selection_id AND users.id = ?", id).Find(&entrees)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
		return nil, result.Error
	}
	return entrees, nil
}

// Maybe create a user (if no errors) and returns the number of inserted records
func CreateEntrees(c context.Context, entrees *[]Entree) error {
	result := db.WithContext(c).Clauses(clause.Returning{}).Create(&entrees)
	return result.Error
}

// Maybe delete a user (if no errors) and returns the number of deleted records
func DeleteEntree(c context.Context, id uuid.UUID) (*int64, error) {
	// Since our models have DeletedAt set, this makes Gorm "soft delete" records on normal delete operations.
	// We can add .Unscoped() prior to the .Delete() call if we want to permanently-delete them.
	result := db.WithContext(c).Delete(&Entree{
		BaseModel: BaseModel{
			ID: id,
		},
	})
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

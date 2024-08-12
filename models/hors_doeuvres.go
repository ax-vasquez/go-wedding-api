package models

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// HorsDouevres table
type HorsDoeuvres struct {
	BaseModel
	OptionName string `json:"option_name" binding:"required"`
}

// Finds all hors doeuvres
func FindHorsDoeuvres(c context.Context) ([]HorsDoeuvres, error) {
	var hors_doeuvres []HorsDoeuvres
	result := db.WithContext(c).Find(&hors_doeuvres)
	return hors_doeuvres, result.Error
}

// Find a single hors doeuvres by ID
func FindHorsDoeuvresById(c context.Context, id uuid.UUID) (*HorsDoeuvres, error) {
	var entree *HorsDoeuvres
	result := db.WithContext(c).Find(&entree, Entree{BaseModel: BaseModel{ID: id}})
	return entree, result.Error
}

// Finds hors doeuvres for the given user
func FindHorsDoeuvresForUser(c context.Context, id uuid.UUID) ([]HorsDoeuvres, error) {
	var hors_doeuvres []HorsDoeuvres
	result := db.WithContext(c).Joins("JOIN users ON hors_doeuvres.id = users.hors_doeuvres_selection_id AND users.id = ?", id).Find(&hors_doeuvres)
	return hors_doeuvres, result.Error
}

// Maybe create a user (if no errors) and returns the number of inserted records
func CreateHorsDoeuvres(c context.Context, hors_douevres *[]HorsDoeuvres) error {
	result := db.WithContext(c).Clauses(clause.Returning{}).Create(&hors_douevres)
	return result.Error
}

// Maybe delete a user (if no errors) and returns the number of deleted records
func DeleteHorsDoeuvres(c context.Context, id uuid.UUID) (*int64, error) {
	// Since our models have DeletedAt set, this makes Gorm "soft delete" records on normal delete operations.
	// We can add .Unscoped() prior to the .Delete() call if we want to permanently-delete them.
	result := db.WithContext(c).Delete(&HorsDoeuvres{}, id)
	return &result.RowsAffected, result.Error
}

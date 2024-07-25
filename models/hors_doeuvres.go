package models

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// HorsDouevres table
type HorsDoeuvres struct {
	BaseModel
	OptionName string `json:"option_name" binding:"required"`
}

// Finds all hors doeuvres
func FindHorsDoeuvres() []HorsDoeuvres {
	var hors_doeuvres []HorsDoeuvres
	result := db.Find(&hors_doeuvres)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
	}
	return hors_doeuvres
}

// Find a single hors doeuvres by ID
func FindHorsDoeuvresById(id uuid.UUID) (*HorsDoeuvres, error) {
	var entree *HorsDoeuvres
	result := db.Find(&entree, Entree{BaseModel: BaseModel{ID: id}})
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
		return nil, result.Error
	}
	return entree, nil
}

// Finds hors doeuvres for the given user
func FindHorsDoeuvresForUser(id uuid.UUID) []HorsDoeuvres {
	var hors_doeuvres []HorsDoeuvres
	result := db.Joins("JOIN users ON hors_doeuvres.id = users.hors_doeuvres_selection_id AND users.id = ?", id).Find(&hors_doeuvres)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
	}
	return hors_doeuvres
}

// Maybe create a user (if no errors) and returns the number of inserted records
func CreateHorsDoeuvres(hors_douevres *[]HorsDoeuvres) (*[]HorsDoeuvres, error) {
	result := db.Clauses(clause.Returning{}).Select("*").Create(&hors_douevres)
	if result.Error != nil {
		return nil, result.Error
	}
	return hors_douevres, nil
}

// Maybe delete a user (if no errors) and returns the number of deleted records
func DeleteHorsDoeuvres(id uuid.UUID) (*int64, error) {
	// Since our models have DeletedAt set, this makes Gorm "soft delete" records on normal delete operations.
	// We can add .Unscoped() prior to the .Delete() call if we want to permanently-delete them.
	result := db.Delete(&HorsDoeuvres{}, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

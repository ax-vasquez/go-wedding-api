package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// HorsDouevres table
type HorsDoeuvres struct {
	gorm.Model
	CreatedAt  time.Time `gorm:"<-:create"`
	OptionName string    `json:"option_name" binding:"required"`
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

// Finds hors doeuvres for the given user
func FindHorsDoeuvresForUser(id uint) []HorsDoeuvres {
	var hors_doeuvres []HorsDoeuvres
	result := db.Joins("JOIN users ON hors_doeuvres.id = users.hors_doeuvres_selection_id AND users.id = ?", id).Find(&hors_doeuvres)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error.Error())
	}
	return hors_doeuvres
}

// Maybe create a user (if no errors) and returns the number of inserted records
func CreateHorsDoeuvres(hors_douevres_opt *HorsDoeuvres) (*int64, error) {
	result := db.Create(&hors_douevres_opt)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

// Maybe delete a user (if no errors) and returns the number of deleted records
func DeleteHorsDoeuvres(id uint) (*int64, error) {
	// Since our models have DeletedAt set, this makes Gorm "soft delete" records on normal delete operations.
	// We can add .Unscoped() prior to the .Delete() call if we want to permanently-delete them.
	result := db.Delete(&HorsDoeuvres{}, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

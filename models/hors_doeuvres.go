package models

import (
	"time"

	"gorm.io/gorm"
)

// HorsDouevres table
type HorsDoeuvres struct {
	gorm.Model
	CreatedAt  time.Time `gorm:"<-:create"`
	OptionName string    `json:"option_name"`
}

// Maybe create a user (if no errors) and returns the number of inserted records
func CreateHorsDoeuvresOption(hors_douevres_opt *HorsDoeuvres) (*int64, error) {
	result := db.Create(&hors_douevres_opt)
	if result.Error != nil {
		return nil, result.Error
	}
	return &result.RowsAffected, nil
}

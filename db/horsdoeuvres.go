package db

import "gorm.io/gorm"

// HorsDouevres table
type HorsDouevres struct {
	gorm.Model
	OptionName string `json:"option_name"`
}

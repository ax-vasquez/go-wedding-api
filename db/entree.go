package db

import "gorm.io/gorm"

// Entree table
type Entree struct {
	gorm.Model
	OptionName string `json:"option_name"`
}

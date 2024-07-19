package models

import (
	"time"

	"gorm.io/gorm"
)

// Entree table
type Entree struct {
	gorm.Model
	CreatedAt  time.Time `gorm:"<-:create"`
	OptionName string    `json:"option_name"`
}

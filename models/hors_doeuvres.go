package models

import (
	"time"

	"gorm.io/gorm"
)

// HorsDouevres table
type HorsDouevres struct {
	gorm.Model
	CreatedAt  time.Time `gorm:"<-:create"`
	OptionName string    `json:"option_name"`
}

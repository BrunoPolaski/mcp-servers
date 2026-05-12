package entities

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	ZipCode    string  `gorm:"size:16;not null"`
	State      string  `gorm:"size:100;not null"`
	City       string  `gorm:"size:100;not null"`
	Street     string  `gorm:"size:255;not null"`
	Number     string  `gorm:"size:16;not null"`
	Complement *string `gorm:"size:255"`
}

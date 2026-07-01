package entities

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	OriginalName string `gorm:"size:255;not null"`
	Name         string `gorm:"size:255;not null"`
	URL          string `gorm:"size:255;not null"`
	MimeType     string `gorm:"size:50;not null"`
}

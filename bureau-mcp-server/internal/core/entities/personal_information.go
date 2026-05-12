package entities

import (
	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

type PersonalInformation struct {
	gorm.Model
	Name      string                   `gorm:"size:255;not null"`
	Phone     valueobjects.PhoneNumber `gorm:"size:15;not null"`
	Document  valueobjects.Document    `gorm:"size:11;uniqueIndex;not null"`
	AddressID *uint
	FileID    *uint
	Address   *Address
	File      *File
}

package entities

import "gorm.io/gorm"

type PersonAddress struct {
	gorm.Model
	PersonalInformationID uint `gorm:"not null;index"`
	AddressID             uint `gorm:"not null;index"`
}

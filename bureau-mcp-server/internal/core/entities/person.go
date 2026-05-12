package entities

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	PersonalInformationID uint `gorm:"not null"`
	PersonalInformation   *PersonalInformation
}

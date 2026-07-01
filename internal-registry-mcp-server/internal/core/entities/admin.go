package entities

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	PersonalInformationID uint `gorm:"not null"`
	PersonalInformation   *PersonalInformation
}

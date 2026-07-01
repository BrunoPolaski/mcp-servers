package entities

import "gorm.io/gorm"

type Analyst struct {
	gorm.Model
	PersonalInformationID uint `gorm:"not null"`
	PersonalInformation   *PersonalInformation
}

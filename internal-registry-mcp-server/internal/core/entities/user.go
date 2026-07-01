package entities

import (
	valueobjects "github.com/BrunoPolaski/internal-registry-mcp-server/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     valueobjects.Email    `gorm:"uniqueIndex;size:255;not null"`
	Password  valueobjects.Password `gorm:"<-:create;size:255;not null"`
	UserType  valueobjects.UserType `gorm:"size:20;not null"`
	PersonID  *uint
	Person    *Person
	AdminID   *uint
	Admin     *Admin
	AnalystID *uint
	Analyst   *Analyst
}

func (u *User) AdditionalInformationRequired() bool {
	return u.Person == nil &&
		u.Admin == nil &&
		u.Analyst == nil
}

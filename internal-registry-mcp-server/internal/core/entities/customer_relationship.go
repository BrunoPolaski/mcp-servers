package entities

import (
	"time"

	"gorm.io/gorm"
)

// CustomerRelationship holds the internal relationship a person has with the
// financial institution (account age, segment, internal scoring).
type CustomerRelationship struct {
	gorm.Model
	PersonID           uint      `gorm:"not null;uniqueIndex"`
	CustomerSince      time.Time `gorm:"not null"`
	RelationshipMonths int       `gorm:"default:0;index"`
	Segment            string    `gorm:"size:50"` // retail, private, business
	Branch             *string   `gorm:"size:100"`
	IsActive           bool      `gorm:"default:true;index"`

	ChurnRisk     *string `gorm:"size:50"` // low, medium, high
	InternalScore *int    `gorm:"index"`
}

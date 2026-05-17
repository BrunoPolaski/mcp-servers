package entities

import (
	"time"

	"gorm.io/gorm"
)

type DataSource struct {
	gorm.Model
	SourceName       string `gorm:"size:255;not null;uniqueIndex"`
	SourceType       string `gorm:"size:100;not null"` // open_finance, receita_federal, internal, partner
	Description      string `gorm:"type:text"`
	IsActive         bool   `gorm:"default:true"`
	LastSyncDate     *time.Time
	ReliabilityScore *int `gorm:"index"` // 0-100
}

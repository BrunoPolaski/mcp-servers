package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
)

type PersonRelationshipDTO struct {
	ID                 uint       `json:"id"`
	CreatedAt          string     `json:"created_at"`
	UpdatedAt          string     `json:"updated_at"`
	PersonID           uint       `json:"person_id"`
	RelatedPersonID    uint       `json:"related_person_id"`
	RelationType       string     `json:"relation_type"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	EndDate            *time.Time `json:"end_date,omitempty"`
	IsActive           bool       `json:"is_active"`
	VerificationStatus string     `json:"verification_status"`
}

func NewPersonRelationshipDTO(entity *entities.PersonRelationship) *PersonRelationshipDTO {
	return &PersonRelationshipDTO{
		ID:                 entity.ID,
		CreatedAt:          entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:          entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:           entity.PersonID,
		RelatedPersonID:    entity.RelatedPersonID,
		RelationType:       entity.RelationType,
		StartDate:          entity.StartDate,
		EndDate:            entity.EndDate,
		IsActive:           entity.IsActive,
		VerificationStatus: entity.VerificationStatus,
	}
}

func (p PersonRelationshipDTO) ToEntity() *entities.PersonRelationship {
	return &entities.PersonRelationship{
		PersonID:           p.PersonID,
		RelatedPersonID:    p.RelatedPersonID,
		RelationType:       p.RelationType,
		StartDate:          p.StartDate,
		EndDate:            p.EndDate,
		IsActive:           p.IsActive,
		VerificationStatus: p.VerificationStatus,
	}
}

type DataSourceDTO struct {
	ID               uint       `json:"id"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	SourceName       string     `json:"source_name"`
	SourceType       string     `json:"source_type"`
	Description      *string    `json:"description,omitempty"`
	IsActive         bool       `json:"is_active"`
	LastSyncDate     *time.Time `json:"last_sync_date,omitempty"`
	ReliabilityScore *int       `json:"reliability_score,omitempty"`
}

func NewDataSourceDTO(entity *entities.DataSource) *DataSourceDTO {
	return &DataSourceDTO{
		ID:               entity.ID,
		CreatedAt:        entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		SourceName:       entity.SourceName,
		SourceType:       entity.SourceType,
		Description:      entity.Description,
		IsActive:         entity.IsActive,
		LastSyncDate:     entity.LastSyncDate,
		ReliabilityScore: entity.ReliabilityScore,
	}
}

func (d DataSourceDTO) ToEntity() *entities.DataSource {
	return &entities.DataSource{
		SourceName:       d.SourceName,
		SourceType:       d.SourceType,
		Description:      d.Description,
		IsActive:         d.IsActive,
		LastSyncDate:     d.LastSyncDate,
		ReliabilityScore: d.ReliabilityScore,
	}
}

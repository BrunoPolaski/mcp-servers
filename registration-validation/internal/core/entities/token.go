package entities

import "time"

type Token struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TokenHash   string     `gorm:"column:token_hash;size:64;not null;uniqueIndex"`
	Description string     `gorm:"column:description"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	ExpiresAt   *time.Time `gorm:"column:expires_at"`
	LastUsedAt  *time.Time `gorm:"column:last_used_at"`
	IsRevoked   bool       `gorm:"column:is_revoked;not null;default:false"`
}

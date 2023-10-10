package model

import (
	"database/sql"
	"time"
)

type AuthenticationToken struct {
	ID          uint32         `json:"id"`
	AdminUserID uint32         `json:"adminUserId"`
	IPAddress   string         `json:"ipAddress"`
	APIToken    sql.NullString `json:"apiToken"`
	ExpiresAt   time.Time      `json:"expiresAt"`
	CreatedAt   sql.NullTime   `json:"createdAt"`
	UpdatedAt   sql.NullTime   `json:"updatedAt"`
}

func (a *AuthenticationToken) TableName() string {
	return "authentication_tokens"
}

package model

import (
	"database/sql"
)

type AdminUserCompany struct {
	UserID    int32         `json:"userId"`
	CompanyID sql.NullInt64 `json:"companyId"`
	CreatedAt sql.NullTime  `json:"createdAt"`
	UpdatedAt sql.NullTime  `json:"updatedAt"`
}

// TableName sets the insert table name for this struct type
func (a *AdminUserCompany) TableName() string {
	return "admin_user_companies"
}

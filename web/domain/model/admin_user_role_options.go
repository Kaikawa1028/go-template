package model

import (
	"database/sql"
)

type AdminUserRoleOption struct {
	UserID            int32         `json:"userId"`
	RoleID            int32         `json:"roleId"`
	DisabledStoreEdit bool          `json:"disabledStoreEdit"`
	CreatedBy         sql.NullInt64 `json:"createdBy"`
	UpdatedBy         sql.NullInt64 `json:"updatedBy"`
	CreatedAt         sql.NullTime  `json:"createdAt"`
	UpdatedAt         sql.NullTime  `json:"updatedAt"`
}

// TableName sets the insert table name for this struct type
func (a *AdminUserRoleOption) TableName() string {
	return "admin_user_role_options"
}

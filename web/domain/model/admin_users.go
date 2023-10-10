package model

import (
	"database/sql"
)

type AdminUser struct {
	ID              uint32         `json:"id"`
	Username        string         `json:"username"`
	Password        string         `json:"password"`
	APIToken        sql.NullString `json:"apiToken"`
	Name            string         `json:"name"`
	Email           sql.NullString `json:"email"`
	Avatar          sql.NullString `json:"avatar"`
	UserID          sql.NullInt64  `json:"userId"`
	Percent         int32          `json:"percent"`
	IndoorPercent   int32          `json:"indoorPercent"`
	Keywords        int32          `json:"keywords"`
	StoreCount      sql.NullInt64  `json:"storeCount"`
	RememberToken   sql.NullString `json:"rememberToken"`
	PayjpCustomerID sql.NullString `json:"payjpCustomerId"`
	FbID            sql.NullInt64  `json:"fbId"`
	FbToken         sql.NullString `json:"fbToken"`
	LoginCount      int32          `json:"loginCount"`
	Exist           sql.NullBool   `json:"exist"`
	CreatedAt       sql.NullTime   `json:"createdAt"`
	UpdatedAt       sql.NullTime   `json:"updatedAt"`
	DeletedAt       sql.NullTime   `json:"deletedAt"`
	CreatedBy       sql.NullInt64  `json:"createdBy"`
	UpdatedBy       sql.NullInt64  `json:"updatedBy"`
	SpecialUser     bool           `json:"specialUser"` // 0: not user special, 1: is user special

}

func (a *AdminUser) TableName() string {
	return "admin_users"
}

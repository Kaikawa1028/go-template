package model

import (
	"database/sql"
)

type Company struct {
	ID              uint32         `json:"id"`
	CompanyID       sql.NullInt64  `json:"companyId"`
	RoleID          bool           `json:"roleId"`
	Credit          bool           `json:"credit"`
	Name            string         `json:"name"`
	PayjpCustomerID sql.NullString `json:"payjpCustomerId"`
	CreatedAt       sql.NullTime   `json:"createdAt"`
	UpdatedAt       sql.NullTime   `json:"updatedAt"`
	PlanID          sql.NullInt64  `json:"planId"`
	CreatedBy       sql.NullInt64  `json:"createdBy"`
	UpdatedBy       sql.NullInt64  `json:"updatedBy"`
	SecretKey       sql.NullString `json:"secretKey"`
	UsedCanlyhp     sql.NullBool   `json:"usedCanlyhp"`
	UsedCanlyhr     bool           `json:"usedCanlyhr"`
	HpAreaSetID     sql.NullInt64  `json:"hpAreaSetId"`
	HpAutoExportCsv bool           `json:"hpAutoExportCsv"`

}

func (c *Company) TableName() string {
	return "companies"
}

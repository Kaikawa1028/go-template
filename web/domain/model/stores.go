package model

import (
	"database/sql"
)

type Store struct {
	ID                uint32         `json:"id"`
	ParentID          int32          `json:"parentId"`
	Title             sql.NullString `json:"title"`
	Email             sql.NullString `json:"email"`
	CardToken         sql.NullString `json:"cardToken"`
	PayjpCustomerID   sql.NullString `json:"payjpCustomerId"`
	Scraping          sql.NullString `json:"scraping"`
	Location          sql.NullInt64  `json:"location"`
	Contract          int32          `json:"contract"`
	Payment           sql.NullInt64  `json:"payment"`
	IndoorCost        int32          `json:"indoorCost"`
	Price             sql.NullInt64  `json:"price"`
	LimitPrice        int32          `json:"limitPrice"`
	AgentRule         bool           `json:"agentRule"`
	LocationID        sql.NullString `json:"locationId"`
	PhoneNumber       sql.NullString `json:"phoneNumber"`
	PostalCode        sql.NullString `json:"postalCode"`
	Available         bool           `json:"available"`
	StartedAt         sql.NullTime   `json:"startedAt"`
	CanceledAt        sql.NullTime   `json:"canceledAt"`
	MeoStartedAt      sql.NullTime   `json:"meoStartedAt"`
	ContinuedAt       sql.NullTime   `json:"continuedAt"`
	CreatedAt         sql.NullTime   `json:"createdAt"`
	UpdatedAt         sql.NullTime   `json:"updatedAt"`
	ScrapingStartTime sql.NullTime   `json:"scrapingStartTime"`
	FbID              sql.NullInt64  `json:"fbId"`
	FbToken           sql.NullString `json:"fbToken"`
	DeletedAt         sql.NullTime   `json:"deletedAt"`
	CreatedBy         sql.NullInt64  `json:"createdBy"`
	UpdatedBy         sql.NullInt64  `json:"updatedBy"`
	Latitude          sql.NullString `json:"latitude"`
	Longitude         sql.NullString `json:"longitude"`
	StoreNumber       sql.NullString `json:"storeNumber"`

	Company  *Company  `gorm:"foreignKey:ParentID"`
}

func (s *Store) TableName() string {
	return "stores"
}

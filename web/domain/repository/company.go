package repository

import (
	gormModel "github.com/Kaikawa1028/go-template/app/domain/model"
	"gorm.io/gorm"
)

type Company interface {
	GetCompanyName(db *gorm.DB, id int) (company *gormModel.Company, err error)
	
}

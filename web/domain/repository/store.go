package repository

import (
	gormModel "github.com/Kaikawa1028/go-template/app/domain/model"
	"gorm.io/gorm"
)

type Store interface {
	ExistStoreRelatedUserIds(db *gorm.DB, storeId uint, userIds []uint32) (bool, error)
	GetStoreByStoreId(db *gorm.DB, storeId uint32) (*gormModel.Store, error)
	GetStoreIdsByCompanyId(db *gorm.DB, companyId uint32) ([]uint32, error)
}

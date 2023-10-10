package repository

import (
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"gorm.io/gorm"
)

type User interface {
	GetByApiToken(db *gorm.DB, apiToken string) (*model.User, error)
	GetByUserId(db *gorm.DB, userId int) (*model.User, error)
	GetAdminUsersRelatedAgencyAdminUser(db *gorm.DB, userId int) ([]*model.User, error)
	UpdateTokenExpiresAt(db *gorm.DB, apiToken string) error
	CanAccessCompany(db *gorm.DB, userIds []uint32, companyId uint32) (bool, error)
	GetAdminUserIdsRelatedAgencyAdminUser(db *gorm.DB, userId int) ([]uint32, error)
}

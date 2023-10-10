//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE

package repository

import (
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"gorm.io/gorm"
)

type AdminUserCompany interface {
	FindByUserIdsAndCompanyIds(db *gorm.DB, userIds []uint32, companyIds []uint32) (records []*model.AdminUserCompany, err error)
}

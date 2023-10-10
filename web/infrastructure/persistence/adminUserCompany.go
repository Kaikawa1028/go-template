package persistence

import (
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"github.com/Kaikawa1028/go-template/app/domain/system"
	"github.com/Kaikawa1028/go-template/app/errors"
	"gorm.io/gorm"
)

type AdminUserCompany struct {
	timer system.ITimer
}

func NewAdminUserCompany(timer system.ITimer) repository.AdminUserCompany {
	return &AdminUserCompany{
		timer: timer,
	}
}

func (c *AdminUserCompany) FindByUserIdsAndCompanyIds(db *gorm.DB, userIds []uint32, companyIds []uint32) (records []*model.AdminUserCompany, err error) {
	err = db.Model(&model.AdminUserCompany{}).
		Select([]string{"user_id", "company_id"}).
		Where("user_id IN (?) AND company_id IN (?)", userIds, companyIds).
		Find(&records).
		Error
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return records, nil
}

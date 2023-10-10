package persistence

import (
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"github.com/Kaikawa1028/go-template/app/errors"
	"gorm.io/gorm"
)

type AdminUserRoleOption struct {
}

func NewAdminUserRoleOption() repository.AdminUserRoleOption {
	return &AdminUserRoleOption{}
}

func (c *AdminUserRoleOption) GetAdminUserRoleOption(db *gorm.DB, userId int) (roleOption *model.AdminUserRoleOption, err error) {
	err = db.Select([]string{"user_id", "role_id", "disabled_store_edit"}).
		Where("user_id = ?", userId).
		First(&roleOption).
		Error
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return roleOption, nil
}

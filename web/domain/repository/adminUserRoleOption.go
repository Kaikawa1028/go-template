package repository

import (
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"gorm.io/gorm"
)

type AdminUserRoleOption interface {
	GetAdminUserRoleOption(db *gorm.DB, userId int) (roleOption *model.AdminUserRoleOption, err error)
}

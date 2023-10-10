package persistence

import (
	"time"

	"github.com/Kaikawa1028/go-template/app/errors/types"

	stdErrors "errors"

	"github.com/Kaikawa1028/go-template/app/config"
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"github.com/Kaikawa1028/go-template/app/errors"
	"gorm.io/gorm"
)

type User struct {
}

func NewUser() repository.User {
	return &User{}
}

func (up *User) GetByApiToken(db *gorm.DB, apiToken string) (*model.User, error) {
	query := `
  SELECT
    users.id,
    users.username,
    users.name,
    users.email,
    roles.slug as role,
    tokens.api_token as token,
    tokens.expires_at
  FROM
    admin_users AS users
    JOIN admin_role_users AS role_users ON users.id = role_users.user_id
    JOIN admin_roles AS roles ON roles.id = role_users.role_id
    JOIN (
      SELECT
        admin_user_id,
        api_token,
        expires_at
      FROM
        authentication_tokens
      WHERE
        api_token = ?
    ) AS tokens ON users.id = tokens.admin_user_id
  WHERE
    deleted_at IS NULL;
    `

	user := &model.User{}

	result := db.Raw(query, apiToken).Scan(user)
	err := result.Error
	if err != nil {
		return nil, errors.Wrap(err)
	}

	//レコードが存在しなかった時
	if result.RowsAffected == 0 {
		err = types.NewNotFoundUserMatchedAuthTokenError()
		return nil, errors.Wrap(err)
	}

	return user, nil
}

func (up *User) GetByUserId(db *gorm.DB, userId int) (*model.User, error) {
	query := `
  SELECT
    users.id,
    users.username,
    users.name,
    COALESCE(users.email, "") AS email,
    roles.slug as role
  FROM
    admin_users AS users
    JOIN admin_role_users AS role_users ON users.id = role_users.user_id
    JOIN admin_roles AS roles ON roles.id = role_users.role_id
  WHERE
    users.id = ?
	AND deleted_at IS NULL;
    `

	user := &model.User{}
	err := db.Raw(query, userId).Scan(user).Error
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return user, nil
}

// 代理店ユーザの配下のユーザを取得します
func (up *User) GetAdminUsersRelatedAgencyAdminUser(db *gorm.DB, userId int) ([]*model.User, error) {
	var query = `
SELECT id, username, name, COALESCE(email, "") AS email
FROM admin_users
WHERE user_id = ?
AND deleted_at IS NULL
`

	var users []*model.User
	err := db.Raw(query, userId).Scan(users).Error
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return users, nil
}

// GetAdminUserIdsRelatedAgencyAdminUser 代理店ユーザの配下のユーザを取得します(Gorm版)
func (up *User) GetAdminUserIdsRelatedAgencyAdminUser(db *gorm.DB, userId int) ([]uint32, error) {
	userIds := []uint32{}

	var users []*model.AdminUser
	err := db.Select([]string{"id"}).
		Where("user_id", userId).
		Where("deleted_at IS NULL").
		Find(&users).
		Error
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return []uint32{}, nil
		}
		return nil, errors.Wrap(err)
	}

	for _, user := range users {
		userIds = append(userIds, user.ID)
	}

	return userIds, nil
}

func (up *User) UpdateTokenExpiresAt(db *gorm.DB, apiToken string) error {
	sessionConfig, err := config.NewSessionConfig()
	if err != nil {
		return errors.Wrap(err)
	}
	currenTime := time.Now()
	expiresAt := currenTime.Add(time.Duration(sessionConfig.Time) * time.Minute)

	result := db.Model(&model.AuthenticationToken{}).
		Where("api_token = ?", apiToken).
		Updates(map[string]interface{}{"expires_at": expiresAt, "updated_at": currenTime})
	err = result.Error

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (up *User) CanAccessCompany(db *gorm.DB, userIds []uint32, companyId uint32) (bool, error) {
	var count int64

	err := db.Model(model.AdminUserCompany{}).
		Where("user_id IN ?", userIds).
		Where("company_id = ?", companyId).
		Count(&count).
		Error

	if err != nil {
		return false, errors.Wrap(err)
	}

	return count > 0, nil
}

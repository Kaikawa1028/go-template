package persistence

import (
	stdErrors "errors"

	gormModel "github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/errors/types"
	"gorm.io/gorm"
)

type StoreStore struct {
	userRepo repository.User
}

func NewStoreStore(
	userRepo repository.User,
) repository.Store {
	return &StoreStore{
		userRepo,
	}
}

func (s *StoreStore) ExistStoreRelatedUserIds(db *gorm.DB, storeId uint, userIds []uint32) (bool, error) {
	query := `
	SELECT COUNT(*) AS count
	FROM stores
	WHERE
		id = ?
		AND (
			parent_id IN (
				SELECT company_id
				FROM admin_user_companies
				WHERE user_id IN ?
			)
			OR
			id IN (
				SELECT store_id
				FROM admin_user_stores
				WHERE admin_user_id IN ?
			)
		)
`

	var result struct {
		Count int `db:"count"`
	}
	err := db.Raw(query, storeId, userIds, userIds).Scan(&result).Error
	if err != nil {
		return false, errors.Wrap(err)
	}

	exist := result.Count >= 1

	return exist, nil
}



func (s *StoreStore) GetStoreByStoreId(db *gorm.DB, storeId uint32) (*gormModel.Store, error) {
	var store *gormModel.Store

	err := db.
		Select([]string{"id", "parent_id", "title", "email"}).
		Where("id = ?", storeId).
		First(&store).
		Error

	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			err = types.NewResourceNotFoundError()
		}
		return nil, errors.Wrap(err)
	}

	return store, nil
}

func (s *StoreStore) GetStoreIdsByCompanyId(db *gorm.DB, companyId uint32) ([]uint32, error) {
	var stores []uint32

	err := db.
		Select([]string{"id"}).
		Where("parent_id = ?", companyId).
		Model(&gormModel.Store{}).
		Find(&stores).
		Error

	if err != nil {
		return nil, errors.Wrap(err)
	}

	return stores, nil
}

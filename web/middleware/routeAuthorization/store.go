package routeAuthorization

import (
	"github.com/labstack/echo/v4"
	domainModel "github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/errors/types"
	"github.com/Kaikawa1028/go-template/app/logger"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Store struct {
	db        *gorm.DB
	storeRepo repository.Store
	userRepo  repository.User
}

func NewStore(
	db *gorm.DB,
	storeRepo repository.Store,
	userRepo repository.User,
) *Store {
	return &Store{
		db,
		storeRepo,
		userRepo,
	}
}

func (m Store) AuthorizeStore(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*domainModel.User)

		// パラメータ取得処理
		// c.Bind()を使わないように注意
		// 1リクエストの間でc.Bind()は1回しか呼べないため、handlerでエラーになります。
		storeId64, err := strconv.ParseUint(c.Param("store-id"), 10, 64)
		if err != nil {
			logger.WarnWithError(c, err, nil)
			return c.NoContent(http.StatusNotFound)
		}
		storeId := uint(storeId64)

		// 権限チェック処理
		canAccess, err := m.canAccess(user, storeId)
		if err != nil {
			logger.Error(c, err, nil)
			return c.NoContent(http.StatusInternalServerError)
		}
		if !canAccess {
			logger.WarnWithError(c, types.NewResourceNotPermittedError(), nil)
			return c.NoContent(http.StatusForbidden)
		}

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}

func (m Store) canAccess(user *domainModel.User, storeId uint) (bool, error) {
	if user.IsAdministrator() {
		return true, nil
	}

	userIds := []uint32{uint32(user.ID)}

	if user.IsAgency() {
		ids, err := m.userRepo.GetAdminUserIdsRelatedAgencyAdminUser(m.db, user.ID)
		if err != nil {
			return false, errors.Wrap(err)
		}
		userIds = append(userIds, ids...)
	}

	exist, err := m.storeRepo.ExistStoreRelatedUserIds(m.db, storeId, userIds)
	if err != nil {
		return false, errors.Wrap(err)
	}

	return exist, nil
}

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

type Company struct {
	db       *gorm.DB
	userRepo repository.User
}

func NewCompany(db *gorm.DB, userRepo repository.User) *Company {
	return &Company{
		db,
		userRepo,
	}
}

func (m Company) AuthorizeCompany(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*domainModel.User)

		// パラメータ取得処理
		// c.Bind()を使わないように注意
		// 1リクエストの間でc.Bind()は1回しか呼べないため、handlerでエラーになります。
		companyId64, err := strconv.ParseUint(c.Param("company-id"), 10, 64)
		if err != nil {
			logger.WarnWithError(c, err, nil)
			return c.NoContent(http.StatusNotFound)
		}
		companyId := uint32(companyId64)

		// 権限チェック処理
		canAccess, err := m.canAccess(user, companyId)
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

func (m Company) canAccess(user *domainModel.User, companyId uint32) (bool, error) {
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

	return m.userRepo.CanAccessCompany(m.db, userIds, companyId)
}

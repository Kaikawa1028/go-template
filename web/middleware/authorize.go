package middleware

import (
	"net/http"

	"github.com/Kaikawa1028/go-template/app/const/app"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/logger"
)

type Authorize struct {
	db                  *gorm.DB
	
}

func NewAuthorize(
	db *gorm.DB,
) *Authorize {
	return &Authorize{
		db,
	}
}

func (m Authorize) Authorize(acceptRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isThrough := false
			if isThrough {
				if err := next(c); err != nil {
					c.Error(err)
				}
				return nil
			}

			user := c.Get("user").(*model.User)

			isAuthorized := false
			if user.Role == app.RoleAdministrator {
				isAuthorized = true
			} else {
				var err error
				isAuthorized, err = m.checkRole(c, acceptRoles, user)
				if err != nil {
					logger.Error(c, err, nil)
					return c.NoContent(http.StatusInternalServerError)
				}
			}
			if !isAuthorized {
				logger.Error(c, errors.New("Authorization failed"), nil)
				return c.NoContent(http.StatusForbidden)
			}

			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}

// checkRole ユーザのロールがAPIの要求するロールを満たしているかどうかをチェックします
// 仕様書： https://www.notion.so/canlyhp/0203fa1332984eb09f0b8069e1955a3c#c60bc615fc804f5eb3953e72fc51a161
func (m Authorize) checkRole(c echo.Context, acceptRoles []string, user *model.User) (bool, error) {
	for _, acceptRole := range acceptRoles {
		if acceptRole == app.RoleWorkerCanEditStore {
			if user.IsWorkerCanEditStore() {
				return true, nil
			}
		} else if acceptRole == user.Role {
			return true, nil
		}
	}

	return false, nil
}


package middleware

import (
	"github.com/Kaikawa1028/go-template/app/errors/types"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Kaikawa1028/go-template/app/const/app"
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"

	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/logger"
	"github.com/Kaikawa1028/go-template/app/usecase"
)

type Authenticate struct {
	authUseCase *usecase.Authenticate
	db          *gorm.DB
	user        repository.User
}

func NewAuthenticate(
	authUseCase *usecase.Authenticate,
	db *gorm.DB,
	user repository.User,
) *Authenticate {
	return &Authenticate{
		authUseCase: authUseCase,
		db:          db,
		user:        user,
	}
}

func (m Authenticate) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isThrough := false
		noAuthPaths := []string{"/v1/healthcheck"}
		for _, path := range noAuthPaths {
			if c.Request().RequestURI == path {
				isThrough = true
				break
			}
		}
		if isThrough {
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}

		var user *model.User
		var err error
		userId := c.Request().Header.Get("x-user-id")
		if os.Getenv("APP_ENV") == app.Local && userId != "" {
			userId, err := strconv.Atoi(userId)
			if err != nil {
				return errors.Wrap(err)
			}
			user, err = m.authUseCase.AuthenticateByUserId(userId)
			if err != nil {
				return errors.Wrap(err)
			}
		} else {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				logger.WarnWithError(c, errors.New("No authentication in request"), nil)
				return c.NoContent(http.StatusUnauthorized)
			}

			user, err = m.authUseCase.Authenticate(token)
			if err != nil {
				if cantBase64DecodeAuthTokenError := types.AsCantBase64DecodeAuthTokenError(err); cantBase64DecodeAuthTokenError != nil {
					logger.WarnWithError(c, err, nil)
					return c.NoContent(http.StatusUnauthorized)
				} else if notFoundUserMatchedAuthTokenError := types.AsNotFoundUserMatchedAuthTokenError(err); notFoundUserMatchedAuthTokenError != nil {
					logger.WarnWithError(c, err, nil)
					return c.NoContent(http.StatusUnauthorized)
				} else {
					return errors.Wrap(err)
				}
			}

			if user == nil {
				logger.WarnWithError(c, errors.New("Not found user who matched token"), nil)
				return c.NoContent(http.StatusUnauthorized)
			}

			if time.Now().After(user.ExpiresAt) {
				logger.WarnWithError(c, errors.New("Expires token"), nil)
				return c.NoContent(http.StatusUnauthorized)
			}

			if err := m.user.UpdateTokenExpiresAt(m.db, user.Token); err != nil {
				return errors.Wrap(err)
			}
		}

		c.Set("user", user)

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

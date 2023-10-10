//go:build integration

package middleware_test

import (
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Kaikawa1028/go-template/app/usecase"
	"gorm.io/gorm"

	"github.com/Kaikawa1028/go-template/app/middleware"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate_ヘッダーに正しいトークンがある場合認証されること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	SetGetUserTestData(db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sample", nil)
	req.Header.Set("Authorization", "tn59HCrmVLGmHPO5FPRt9yLHfKk%2Fgxt%2F%2FU03i%2FuBVG%2BEwbsexDcTvZVwztHdy1GE7kjV1T1rDAePGFJOOW%2F3dVv5%2BdBi4233NUWvNiZ%2Fc4JAlIgwtTKDaQ0i%2B3eUzf%2Bj")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	m := middleware.NewAuthenticate(
		usecase.NewAuthenticate(
			cont.DI.OrmDB,
			cont.DI.User,
			cont.DI.AdminUserRoleOption,
		),
		cont.DI.OrmDB,
		cont.DI.User,
	)
	h := m.Authenticate(func(c echo.Context) error {
		return nil
	})
	err := h(c)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	user := c.Get("user").(*model.User)
	assert.Equal(t, 1001, user.ID)
	assert.Equal(t, "testadmin", user.Username)
	assert.Equal(t, "testadmin", user.Name)
	assert.Equal(t, "admin1@example.com", user.Email)
	assert.Equal(t, "administrator", user.Role)
	assert.Equal(t, "19502240fffc01e84fec5b850e9bf124bbc53a55269003f608934141ec3ab216", user.Token)
	assert.Nil(t, user.RoleOption)
}

func TestAuthenticate_ヘッダーに間違ってるトークンがある場合認証されないこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	SetGetUserTestData(db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sample", nil)
	req.Header.Set("Authorization", "en59HCrmVLGmHPO5FPRt9yLHfKk%2Fgxt%2F%2FU03i%2FuBVG%2BEwbsexDcTvZVwztHdy1GE7kjV1T1rDAePGFJOOW%2F3dVv5%2BdBi4233NUWvNiZ%2Fc4JAlIgwtTKDaQ0i%2B3eUzf%2Bj")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	m := middleware.NewAuthenticate(
		usecase.NewAuthenticate(
			cont.DI.OrmDB,
			cont.DI.User,
			cont.DI.AdminUserRoleOption,
		),
		cont.DI.OrmDB,
		cont.DI.User,
	)
	h := m.Authenticate(func(c echo.Context) error {
		return nil
	})

	err := h(c)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthenticate_ヘッダーにBase64形式ではないトークンがある場合認証されないこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	SetGetUserTestData(db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sample", nil)
	req.Header.Set("Authorization", "en59HCrmVLGmHPO5FPRt9yL")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	m := middleware.NewAuthenticate(
		usecase.NewAuthenticate(
			cont.DI.OrmDB,
			cont.DI.User,
			cont.DI.AdminUserRoleOption,
		),
		cont.DI.OrmDB,
		cont.DI.User,
	)
	h := m.Authenticate(func(c echo.Context) error {
		return nil
	})

	err := h(c)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthenticate_ヘッダーにトークンがないエラーが返却されること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sample", nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	m := middleware.NewAuthenticate(
		usecase.NewAuthenticate(
			cont.DI.OrmDB,
			cont.DI.User,
			cont.DI.AdminUserRoleOption,
		),
		cont.DI.OrmDB,
		cont.DI.User,
	)
	h := m.Authenticate(func(c echo.Context) error {
		return nil
	})

	err := h(c)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthenticate_トークンの有効期限が切れてた場合認証されないこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	SetGetUserTestData(db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sample", nil)
	req.Header.Set("Authorization", "%2FBQZ7TTn8P0pbksaL%2FifXMWFJitgaN0yzh2pxJ32YNAapg0jcd%2BWdri8u6%2FdlUL00vbjfvfJBdZJzN2%2FYD3Hf9srj67zFHGNdGUrc76%2Ba%2FKv4ZvPOJshWny2Rqa3QBEh")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	m := middleware.NewAuthenticate(
		usecase.NewAuthenticate(
			cont.DI.OrmDB,
			cont.DI.User,
			cont.DI.AdminUserRoleOption,
		),
		cont.DI.OrmDB,
		cont.DI.User,
	)
	h := m.Authenticate(func(c echo.Context) error {
		return nil
	})

	err := h(c)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthenticate_認証後トークンの有効期限が延長されること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	SetGetUserTestData(db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sample", nil)
	req.Header.Set("Authorization", "tn59HCrmVLGmHPO5FPRt9yLHfKk%2Fgxt%2F%2FU03i%2FuBVG%2BEwbsexDcTvZVwztHdy1GE7kjV1T1rDAePGFJOOW%2F3dVv5%2BdBi4233NUWvNiZ%2Fc4JAlIgwtTKDaQ0i%2B3eUzf%2Bj")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	m := middleware.NewAuthenticate(
		usecase.NewAuthenticate(
			cont.DI.OrmDB,
			cont.DI.User,
			cont.DI.AdminUserRoleOption,
		),
		cont.DI.OrmDB,
		cont.DI.User,
	)
	h := m.Authenticate(func(c echo.Context) error {
		return nil
	})

	assert.Equal(t, nil, h(c))

	var Token struct {
		ExpiresAt time.Time `db:"expires_at"`
	}
	sql := "SELECT expires_at FROM authentication_tokens WHERE admin_user_id = 1001"
	err := db.Raw(sql).Scan(&Token).Error
	if err != nil {
		panic(err)
	}

	comparisonTime := time.Now().Add(90 * time.Minute)

	assert.Equal(t, true, Token.ExpiresAt.After(comparisonTime))
}

func TestAuthenticate_admin_user_role_optionsテーブルにレコードが存在する場合ロールオプションが取得できること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	expiresAt1 := time.Now().Add(60 * time.Minute)
	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "admin_users", []map[string]interface{}{
			{"id": 1001, "username": "worker", "name": "worker", "password": "", "email": "worker@example.com"},
		})
		mustInsert(db, "admin_role_users", []map[string]interface{}{
			{"role_id": 13, "user_id": 1001},
		})
		mustInsert(db, "admin_user_role_options", []map[string]interface{}{
			{"role_id": 13, "user_id": 1001, "disabled_store_edit": true},
		})
		mustInsert(db, "authentication_tokens", []map[string]interface{}{
			{"admin_user_id": 1001, "ip_address": "172.000.000.001", "api_token": "19502240fffc01e84fec5b850e9bf124bbc53a55269003f608934141ec3ab216", "expires_at": expiresAt1},
		})
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sample", nil)
	req.Header.Set("Authorization", "tn59HCrmVLGmHPO5FPRt9yLHfKk%2Fgxt%2F%2FU03i%2FuBVG%2BEwbsexDcTvZVwztHdy1GE7kjV1T1rDAePGFJOOW%2F3dVv5%2BdBi4233NUWvNiZ%2Fc4JAlIgwtTKDaQ0i%2B3eUzf%2Bj")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	m := middleware.NewAuthenticate(
		usecase.NewAuthenticate(
			cont.DI.OrmDB,
			cont.DI.User,
			cont.DI.AdminUserRoleOption,
		),
		cont.DI.OrmDB,
		cont.DI.User,
	)
	h := m.Authenticate(func(c echo.Context) error {
		return nil
	})

	assert.Equal(t, nil, h(c))

	user := c.Get("user").(*model.User)
	assert.Equal(t, 1001, user.ID)
	assert.Equal(t, "worker", user.Username)
	assert.Equal(t, "worker", user.Name)
	assert.Equal(t, "worker@example.com", user.Email)
	assert.Equal(t, "worker", user.Role)
	assert.Equal(t, "19502240fffc01e84fec5b850e9bf124bbc53a55269003f608934141ec3ab216", user.Token)
	assert.Equal(t, true, user.RoleOption.DisabledStoreEdit)
}

func SetGetUserTestData(db *gorm.DB) {
	expiresAt1 := time.Now().Add(60 * time.Minute)
	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "admin_users", []map[string]interface{}{
			{"id": 1001, "username": "testadmin", "name": "testadmin", "password": "", "email": "admin1@example.com"},
			{"id": 1002, "username": "testmanager", "name": "testmanager", "password": "", "email": "manager1@example.com"},
			{"id": 1003, "username": "testviewer", "name": "testviewer", "password": "", "email": "viewer1@example.com"},
		})

		mustInsert(db, "admin_role_users", []map[string]interface{}{
			{"role_id": 1, "user_id": 1001},
			{"role_id": 1, "user_id": 1002},
		})
		mustInsert(db, "authentication_tokens", []map[string]interface{}{
			{"admin_user_id": 1001, "ip_address": "172.000.000.001", "api_token": "19502240fffc01e84fec5b850e9bf124bbc53a55269003f608934141ec3ab216", "expires_at": expiresAt1},
			{"admin_user_id": 1002, "ip_address": "172.000.000.001", "api_token": "bed9d04d6cf9ef3a87b997dee9d19e2e77e82cbd954e805cbdc2e8d850b9c05a", "expires_at": "2020-02-16 16:34:10"},
		})
	})
}

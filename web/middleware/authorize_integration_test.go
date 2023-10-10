//go:build integration

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kaikawa1028/go-template/app/const/app"
	"gorm.io/gorm"

	"github.com/Kaikawa1028/go-template/app/middleware"

	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestAuthorize_アドミンの場合認可が通ること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "administrator",
		RoleOption: &model.AdminUserRoleOption{
			DisabledStoreEdit: false,
		},
	})
	m := middleware.NewAuthorize(cont.DI.OrmDB)
	h := m.Authorize([]string{"agency", "manager", "editor", "worker", "viewer"})(func(c echo.Context) error {
		return nil
	})

	assert.Equal(t, nil, h(c))
}

func TestAuthorize_管理者権限のユーザの場合認可が通ること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "manager",
		RoleOption: &model.AdminUserRoleOption{
			DisabledStoreEdit: false,
		},
	})
	m := middleware.NewAuthorize(cont.DI.OrmDB)
	h := m.Authorize([]string{"agency", "manager", "editor", "worker"})(func(c echo.Context) error {
		return nil
	})

	assert.Equal(t, nil, h(c))
}


func TestAuthorize_管理者_システム管理設定可能_を許可しているかつユーザが管理者ではない場合認可が通らないこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "cms_authority_settings", []map[string]interface{}{
			{
				"id":                     1,
				"company_id":             20,
				"manager_format_setting": true,
				"manager_seo_setting":    true,
				"manager_top_setting":    true,
			},
		})
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "viewer",
		RoleOption: &model.AdminUserRoleOption{
			DisabledStoreEdit: false,
		},
	})
	m := middleware.NewAuthorize(cont.DI.OrmDB)
	h := m.Authorize([]string{app.RoleManagerCanFormatSetting})(func(c echo.Context) error {
		return nil
	})

	h(c)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}



func TestAuthorize_閲覧者権限のユーザの場合認可が通らないこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "viewer",
		RoleOption: &model.AdminUserRoleOption{
			DisabledStoreEdit: false,
		},
	})
	m := middleware.NewAuthorize(cont.DI.OrmDB)
	h := m.Authorize([]string{"agency", "manager", "editor", "worker"})(func(c echo.Context) error {
		return nil
	})

	h(c)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthorize_作業者店舗編集可能のみを許可している場合店舗編集が許可されていない作業者が拒否されること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "worker",
		RoleOption: &model.AdminUserRoleOption{
			DisabledStoreEdit: true,
		},
	})
	m := middleware.NewAuthorize(cont.DI.OrmDB)
	h := m.Authorize([]string{"worker:can_edit_store"})(func(c echo.Context) error {
		return nil
	})

	h(c)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthorize_作業者店舗編集可能のみを許可している場合店舗編集が許可されている作業者が許可されること1(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "worker",
		RoleOption: &model.AdminUserRoleOption{
			DisabledStoreEdit: false,
		},
	})
	m := middleware.NewAuthorize(cont.DI.OrmDB)
	h := m.Authorize([]string{"worker:can_edit_store"})(func(c echo.Context) error {
		return nil
	})

	h(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthorize_作業者店舗編集可能のみを許可している場合店舗編集が許可されている作業者が許可されること2(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:         10000,
		Username:   "test",
		Name:       "テストユーザ",
		Email:      "test@test.test",
		Role:       "worker",
		RoleOption: nil, // admin_user_role_optionが存在しないケース
	})
	m := middleware.NewAuthorize(cont.DI.OrmDB)
	h := m.Authorize([]string{"worker:can_edit_store"})(func(c echo.Context) error {
		return nil
	})

	h(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthorize_作業者のみを許可している場合店舗編集が許可されていない作業者でも許可されること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "worker",
		RoleOption: &model.AdminUserRoleOption{
			DisabledStoreEdit: true,
		},
	})
	m := middleware.NewAuthorize(cont.DI.OrmDB)
	h := m.Authorize([]string{"worker"})(func(c echo.Context) error {
		return nil
	})

	h(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

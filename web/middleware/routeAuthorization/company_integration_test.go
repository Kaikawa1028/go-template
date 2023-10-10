//go:build integration

package routeAuthorization_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/middleware/routeAuthorization"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func TestAuthorizeCompany_閲覧者ロールの場合所属している会社にアクセスできること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "admin_user_companies", []map[string]interface{}{
			{"user_id": 30, "company_id": 20},
		})
	})

	user := &model.User{
		ID:       30,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "viewer",
	}
	rec, err := callApiAuthorizeCompany(cont, user, "20")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthorizeCompany_閲覧者ロールの場合所属していない会社にアクセスできないこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "admin_user_companies", []map[string]interface{}{
			{"user_id": 30, "company_id": 20},
		})
	})

	user := &model.User{
		ID:       30,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "viewer",
	}
	rec, err := callApiAuthorizeCompany(cont, user, "21")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthorizeCompany_代理店ロールの場合配下のアカウントが所属している会社にアクセスできること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "admin_user_companies", []map[string]interface{}{
			{"user_id": 40, "company_id": 20},
		})
		mustInsert(db, "admin_users", []map[string]interface{}{
			{"id": 40, "user_id": 30, "username": "DUMMY", "password": "DUMMY", "name": "DUMMY"},
		})
	})

	user := &model.User{
		ID:       30,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "agency",
	}
	rec, err := callApiAuthorizeCompany(cont, user, "20")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthorizeCompany_代理店ロールの場合配下のアカウントが所属していない会社にアクセスできないこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "admin_user_companies", []map[string]interface{}{
			{"user_id": 30, "company_id": 20},
		})
		mustInsert(db, "admin_users", []map[string]interface{}{
			{"id": 40, "user_id": 30, "username": "DUMMY", "password": "DUMMY", "name": "DUMMY"},
		})
	})

	user := &model.User{
		ID:       30,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "agency",
	}
	rec, err := callApiAuthorizeCompany(cont, user, "21")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthorizeCompany_アドミンロールの場合所属していない会社にアクセスできること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "admin_user_companies", []map[string]interface{}{
			{"user_id": 30, "company_id": 20},
		})
	})

	user := &model.User{
		ID:       30,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "administrator",
	}
	rec, err := callApiAuthorizeCompany(cont, user, "21")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
}

func callApiAuthorizeCompany(cont *TestCaseContainer, user *model.User, companyId string) (*httptest.ResponseRecorder, error) {
	e := cont.DI.Server.Echo
	req := httptest.NewRequest(http.MethodGet, "/v1/companies/:company-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/companies/:company-id")
	c.SetParamNames("company-id")
	c.SetParamValues(companyId)
	c.Set("user", user)
	m := routeAuthorization.NewCompany(cont.DI.OrmDB, cont.DI.User)
	h := m.AuthorizeCompany(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	err := h(c)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

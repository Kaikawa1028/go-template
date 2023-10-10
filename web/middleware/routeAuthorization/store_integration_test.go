//go:build integration

package routeAuthorization_test

import (
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/middleware/routeAuthorization"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAuthorizeStore_閲覧者ロールの場合所属している会社に紐付いている店舗にアクセスできること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "stores", []map[string]interface{}{
			{"id": 10, "parent_id": 20},
		})
		mustInsert(db, "companies", []map[string]interface{}{
			{"id": 20, "name": "会社1"},
		})
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
	rec, err := callApiAuthorizeStore(cont, user, "10")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthorizeStore_閲覧者ロールの場合所属していない会社に紐付いている店舗にアクセスできないこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		// 店舗が所属していない会社に紐付いている
		mustInsert(db, "stores", []map[string]interface{}{
			{"id": 10, "parent_id": 99},
		})
		mustInsert(db, "companies", []map[string]interface{}{
			{"id": 20, "name": "会社1"},
		})
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
	rec, err := callApiAuthorizeStore(cont, user, "10")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthorizeStore_閲覧者ロールかつアカウントに店舗が直接紐付けされている場合店舗にアクセスできること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "stores", []map[string]interface{}{
			{"id": 10, "parent_id": 99},
		})
		mustInsert(db, "admin_user_stores", []map[string]interface{}{
			{"admin_user_id": 30, "store_id": 10},
		})
	})

	user := &model.User{
		ID:       30,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "viewer",
	}
	rec, err := callApiAuthorizeStore(cont, user, "10")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthorizeStore_代理店ロールかつ配下のアカウントが所属している会社に紐付いている店舗にアクセスできること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "stores", []map[string]interface{}{
			{"id": 10, "parent_id": 20},
		})
		mustInsert(db, "companies", []map[string]interface{}{
			{"id": 20, "name": "会社1"},
		})
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
	rec, err := callApiAuthorizeStore(cont, user, "10")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthorizeStore_代理店ロールかつ配下のアカウントに店舗が直接紐付けされている場合店舗にアクセスできること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "stores", []map[string]interface{}{
			{"id": 10, "parent_id": 20},
		})
		mustInsert(db, "admin_user_stores", []map[string]interface{}{
			{"admin_user_id": 30, "store_id": 10},
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
	rec, err := callApiAuthorizeStore(cont, user, "10")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthorizeStore_代理店ロールかつ配下のアカウントもアクセス権を持たない場合店舗にアクセスできないこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "stores", []map[string]interface{}{
			{"id": 10, "parent_id": 20},
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
	rec, err := callApiAuthorizeStore(cont, user, "10")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthorizeStore_管理者ロールの場合所属していない会社に紐付いている店舗にアクセスできること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		// 店舗が所属していない会社に紐付いている
		mustInsert(db, "stores", []map[string]interface{}{
			{"id": 10, "parent_id": 99},
		})
		mustInsert(db, "companies", []map[string]interface{}{
			{"id": 20, "name": "会社1"},
		})
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
	rec, err := callApiAuthorizeStore(cont, user, "10")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
}

func callApiAuthorizeStore(cont *TestCaseContainer, user *model.User, storeId string) (*httptest.ResponseRecorder, error) {
	e := cont.DI.Server.Echo
	req := httptest.NewRequest(http.MethodGet, "/v1/stores/:store-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/stores/:store-id")
	c.SetParamNames("store-id")
	c.SetParamValues(storeId)
	c.Set("user", user)
	m := routeAuthorization.NewStore(cont.DI.OrmDB, cont.DI.Store, cont.DI.User)
	h := m.AuthorizeStore(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	err := h(c)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

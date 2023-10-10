//go:build integration

package handler_test

import (
	"encoding/json"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/errors/types"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kaikawa1028/go-template/app/usecase"
	"gorm.io/gorm"

	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/interfaces/handler"
	"github.com/stretchr/testify/assert"
)

func TestGetCompanyName_代理店ロールで会社名を取得できること(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "companies", []map[string]interface{}{
			{"id": 1111, "name": "会社1"},
		})
	})

	//
	// Execute
	//

	e := cont.DI.Server.Echo
	req := httptest.NewRequest(http.MethodGet, "/v1/companies/:company-id/name", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/companies/:company-id/name")
	c.SetParamNames("company-id")
	c.SetParamValues("1111")
	c.Set("user", &model.User{
		ID:       1,
		Username: "testAdmin",
		Name:     "Adminテストユーザ",
		Email:    "admin1@example.com",
		Role:     "administrator",
	})
	h := handler.NewGetCompanyName(
		usecase.NewGetCompanyName(
			cont.DI.OrmDB,
			cont.DI.User,
			cont.DI.Company,
		))
	err := h.GetCompanyName(c)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Assert
	//
	assert.Equal(t, http.StatusOK, rec.Code)

	var response handler.GetCompanyNameResponse
	b, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("Could not read the response body")
	}
	if err := json.Unmarshal(b, &response); err != nil {
		t.Fatalf("Could not unmarshal the response body")
	}

	assert.Equal(t, "会社1", response.Name)
}

func TestGetCompanyName_代理店ロールで存在しない会社IDが指定された場合はResourceNotFoundErrorを返却できること(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)
	db := cont.DI.OrmDB

	prepareTestData(db, func(db *gorm.DB) {
		mustInsert(db, "companies", []map[string]interface{}{
			{"id": 1111, "name": "会社1"},
		})
	})

	//
	// Execute
	//

	e := cont.DI.Server.Echo
	req := httptest.NewRequest(http.MethodGet, "/v1/companies/:company-id/name", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/companies/:company-id/name")
	c.SetParamNames("company-id")
	c.SetParamValues("9999")
	c.Set("user", &model.User{
		ID:       1,
		Username: "testAdmin",
		Name:     "Adminテストユーザ",
		Email:    "admin1@example.com",
		Role:     "administrator",
	})
	h := handler.NewGetCompanyName(
		usecase.NewGetCompanyName(
			cont.DI.OrmDB,
			cont.DI.User,
			cont.DI.Company,
		))
	err := h.GetCompanyName(c)

	//
	// Assert
	//
	assert.ErrorAs(t, errors.GetRootError(err), &types.ResourceNotFoundError{})
}

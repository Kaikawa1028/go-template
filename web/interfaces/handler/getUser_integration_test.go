//go:build integration

package handler_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/interfaces/handler"
	"github.com/stretchr/testify/assert"
)

func TestGetUser_ログインユーザを取得できること(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	user := &model.User{
		ID:       3,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "administrator",
		RoleOption: &model.AdminUserRoleOption{
			DisabledStoreEdit: true,
		},
	}
	rec, err := callApiGetUser(cont, user)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Assert
	//
	assert.Equal(t, http.StatusOK, rec.Code)

	var userRes handler.GetUserResponseUser
	b, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &userRes); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, userRes.Name, "テストユーザ")
	assert.Equal(t, userRes.Role, "administrator")
	assert.Equal(t, userRes.RoleOption.DisabledStoreEdit, true)
}

func TestGetUser_ロールオプションがnilの場合初期値を返すこと(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := cont.DI.Server.Echo
	req := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:         3,
		Username:   "test",
		Name:       "テストユーザ",
		Email:      "test@test.test",
		Role:       "administrator",
		RoleOption: nil,
	})
	h := handler.NewGetUser()
	err := h.GetUser(c)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Assert
	//
	assert.Equal(t, http.StatusOK, rec.Code)

	var userRes handler.GetUserResponseUser
	b, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &userRes); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, userRes.RoleOption.DisabledStoreEdit, false)
}

func callApiGetUser(cont *TestCaseContainer, user *model.User) (*httptest.ResponseRecorder, error) {
	e := cont.DI.Server.Echo
	req := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)
	h := handler.NewGetUser()
	err := h.GetUser(c)
	return rec, err
}

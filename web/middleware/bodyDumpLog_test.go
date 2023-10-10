//go:build integration

package middleware_test

import (
	"encoding/json"
	"github.com/Kaikawa1028/go-template/app/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestBodyDumpLog_LOG_LEVELがdebugの場合BodyDumpログが出力されること(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	req.Header.Set("Referer", "https://example.com/")
	req.Header.Set("Test-Header-Key1", "Test-Header-Value-1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "administrator",
	})

	t.Setenv("LOG_LEVEL", "debug")

	// JSONを返すモックhandlerを作成する
	respBody := struct {
		TestKey1 string `json:"testKey1"`
	}{
		TestKey1: "testValue1",
	}
	mockHandler := func(c echo.Context) error { return c.JSON(http.StatusOK, respBody) }

	//
	// Execute
	//
	m := middleware.NewBodyDumpLog()
	err := m.BodyDumpLog()(mockHandler)(c)

	//
	// Assert
	//
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(loggerHook.Entries))

	logStr, err := loggerHook.LastEntry().String()
	if err != nil {
		t.Fatal(err)
	}

	var log map[string]interface{}
	err = json.Unmarshal([]byte(logStr), &log)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "[Body Dump]", log["message"])
	assert.Equal(t, "debug", log["level"])
	assert.NotEmpty(t, log["function"])
	assert.NotEmpty(t, log["file"])
	assert.NotEmpty(t, log["line"])
	assert.NotEmpty(t, log["host"])
	headers := log["header"].(map[string]interface{})
	assert.Len(t, headers, 2)
	assert.Equal(t, "Test-Header-Value-1", headers["Test-Header-Key1"].([]interface{})[0])
	output := log["output"].(map[string]interface{})
	assert.Equal(t, "testValue1", output["testKey1"])
}

func TestBodyDumpLog_LOG_LEVELがdebug以外の場合BodyDumpログが出力されないこと(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	req.Header.Set("Referer", "https://example.com/")
	req.Header.Set("Test-Header-Key1", "Test-Header-Value-1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "administrator",
	})

	t.Setenv("LOG_LEVEL", "info")

	// JSONを返すモックhandlerを作成する
	respBody := struct {
		TestKey1 string `json:"testKey1"`
	}{
		TestKey1: "testValue1",
	}
	mockHandler := func(c echo.Context) error { return c.JSON(http.StatusOK, respBody) }

	//
	// Execute
	//
	m := middleware.NewBodyDumpLog()
	err := m.BodyDumpLog()(mockHandler)(c)

	//
	// Assert
	//
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 0, len(loggerHook.Entries))
}

package logger_test

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestInfo_Infoログが出力される事(t *testing.T) {
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

	//
	// Execute
	//
	logger.Info(c, "テストメッセージ", map[string]interface{}{
		"testKey1": "testValue1",
	})

	//
	// Assert
	//
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

	assert.Equal(t, "テストメッセージ", log["message"])
	assert.Equal(t, "info", log["level"])
	assert.Equal(t, "testValue1", log["params"].(map[string]interface{})["testKey1"])
	assert.NotEmpty(t, log["function"])
	assert.NotEmpty(t, log["file"])
	assert.NotEmpty(t, log["line"])
	assert.NotEmpty(t, log["host"])
	assert.Equal(t, "/test-path", log["uri"])
	assert.Equal(t, "192.0.2.1", log["ip"])
	assert.Equal(t, "GET", log["http_method"])
	assert.NotEmpty(t, log["server"])
	assert.Equal(t, "https://example.com/", log["referrer"])
	assert.Equal(t, "test", log["environment"])
	assert.Nil(t, log["input"])
	headers := log["header"].(map[string]interface{})
	assert.Len(t, headers, 2)
	assert.Equal(t, "Test-Header-Value-1", headers["Test-Header-Key1"].([]interface{})[0])
	assert.Equal(t, 10000, int(log["userId"].(float64))) // floatとしてパースされてしまうのでintに直す
}

func TestError_Errorログが出力される事(t *testing.T) {
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

	//
	// Execute
	//
	logger.Error(c, fmt.Errorf("テストエラーメッセージ"), map[string]interface{}{
		"testKey1": "testValue1",
	})

	//
	// Assert
	//
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

	assert.Equal(t, "テストエラーメッセージ", log["message"])
	assert.Equal(t, "error", log["level"])
	assert.Equal(t, "testValue1", log["params"].(map[string]interface{})["testKey1"])
	assert.NotEmpty(t, log["function"])
	assert.NotEmpty(t, log["file"])
	assert.NotEmpty(t, log["line"])
	assert.NotEmpty(t, log["host"])
	assert.Equal(t, "/test-path", log["uri"])
	assert.Equal(t, "192.0.2.1", log["ip"])
	assert.Equal(t, "GET", log["http_method"])
	assert.NotEmpty(t, log["server"])
	assert.Equal(t, "https://example.com/", log["referrer"])
	assert.Equal(t, "test", log["environment"])
	assert.Nil(t, log["input"])
	headers := log["header"].(map[string]interface{})
	assert.Len(t, headers, 2)
	assert.Equal(t, "Test-Header-Value-1", headers["Test-Header-Key1"].([]interface{})[0])
	assert.Equal(t, 10000, int(log["userId"].(float64))) // floatとしてパースされてしまうのでintに直す
	assert.Equal(t, false, log["panic"])
}

func TestError_AppErrorを渡した場合スタックトレースを含むログが出力される事(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	reqBody := "{ testKey1: 'testValue1' }"
	req := httptest.NewRequest(http.MethodGet, "/test-path", strings.NewReader(reqBody))
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

	//
	// Execute
	//
	logger.Error(c, errors.New("テストエラーメッセージ"), map[string]interface{}{
		"testKey1": "testValue1",
	})

	//
	// Assert
	//
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

	assert.Regexp(t, "テストエラーメッセージ", log["message"])
	assert.Regexp(t, "logger_test.go:[0-9]+", log["message"])
	assert.Equal(t, "error", log["level"])
	assert.Equal(t, "testValue1", log["params"].(map[string]interface{})["testKey1"])
	assert.NotEmpty(t, log["function"])
	assert.NotEmpty(t, log["file"])
	assert.NotEmpty(t, log["line"])
	assert.NotEmpty(t, log["host"])
	assert.Equal(t, "/test-path", log["uri"])
	assert.Equal(t, "192.0.2.1", log["ip"])
	assert.Equal(t, "GET", log["http_method"])
	assert.NotEmpty(t, log["server"])
	assert.Equal(t, "https://example.com/", log["referrer"])
	assert.Equal(t, "test", log["environment"])
	assert.Nil(t, log["input"])
	headers := log["header"].(map[string]interface{})
	assert.Len(t, headers, 2)
	assert.Equal(t, "Test-Header-Value-1", headers["Test-Header-Key1"].([]interface{})[0])
	assert.Equal(t, 10000, int(log["userId"].(float64))) // floatとしてパースされてしまうのでintに直す
	assert.Equal(t, false, log["panic"])
}

func TestWarnWithError_AppErrorを渡した場合スタックトレースを含むWarningログが出力される事(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	reqBody := "{ testKey1: 'testValue1' }"
	req := httptest.NewRequest(http.MethodGet, "/test-path", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//
	// Execute
	//
	logger.WarnWithError(c, errors.New("テストWarnログ"), map[string]interface{}{
		"testKey1": "testValue1",
	})

	//
	// Assert
	//
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

	assert.Regexp(t, "テストWarnログ", log["message"])
	assert.Regexp(t, "logger_test.go:[0-9]+", log["message"])
	assert.Equal(t, "warning", log["level"])
}

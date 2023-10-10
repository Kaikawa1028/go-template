//go:build integration

package middleware_test

import (
	"github.com/Kaikawa1028/go-template/app/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetupResponseHeader_リクエストヘッダーのXAmznTraceIdがレスポンスヘッダーにコピーされる事(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	req.Header.Set("X-Amzn-Trace-Id", "dummy-trace-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	mockHandler := func(c echo.Context) error { return c.NoContent(http.StatusOK) }

	//
	// Execute
	//
	m := middleware.NewSetupResponseHeader()
	err := m.SetupResponseHeader(mockHandler)(c)

	//
	// Assert
	//
	if err != nil {
		t.Error(err)
	}

	traceId := rec.Header().Get("X-Amzn-Trace-Id")
	assert.Equal(t, traceId, "dummy-trace-id")
}

func TestSetupResponseHeader_リクエストヘッダーにXAmznTraceIdが存在しない場合はレスポンスヘッダーにコピーされない事(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	mockHandler := func(c echo.Context) error { return c.NoContent(http.StatusOK) }

	//
	// Execute
	//
	m := middleware.NewSetupResponseHeader()
	err := m.SetupResponseHeader(mockHandler)(c)

	//
	// Assert
	//
	if err != nil {
		t.Error(err)
	}

	traceId := rec.Header().Get("X-Amzn-Trace-Id")
	assert.Empty(t, traceId)
}

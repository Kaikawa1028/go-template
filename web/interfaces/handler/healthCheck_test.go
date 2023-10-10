package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/Kaikawa1028/go-template/app/interfaces/handler"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewHealthCheck()

	if assert.NoError(t, h.HealthCheck(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

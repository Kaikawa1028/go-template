package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthCheck struct {
}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

func (h HealthCheck) HealthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

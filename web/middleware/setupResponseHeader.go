package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/const/http"
)

type SetupResponseHeader struct {
}

func NewSetupResponseHeader() *SetupResponseHeader {
	return &SetupResponseHeader{}
}

func (m SetupResponseHeader) SetupResponseHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		setAmazonTraceId(c)
		return next(c)
	}
}

func setAmazonTraceId(c echo.Context) {
	traceId := c.Request().Header.Get(http.HeaderAmazonTraceId)
	if traceId != "" {
		c.Response().Header().Set(http.HeaderAmazonTraceId, traceId)
	}
}

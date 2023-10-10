package middleware

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/Kaikawa1028/go-template/app/logger"
	"os"
	"strings"
)

// BodyDumpLog レスポンスBodyをログ出力するミドルウェア
type BodyDumpLog struct {
}

func NewBodyDumpLog() *BodyDumpLog {
	return &BodyDumpLog{}
}

func (m BodyDumpLog) BodyDumpLog() echo.MiddlewareFunc {
	return echoMiddleware.BodyDumpWithConfig(echoMiddleware.BodyDumpConfig{
		Skipper: func(c echo.Context) bool {
			logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
			switch logLevel {
			case "debug", "trace":
				return false
			default:
				return true
			}
		},

		Handler: func(c echo.Context, reqBodyBytes []byte, resBodyBytes []byte) {
			var resBody map[string]interface{}
			resContentType := c.Response().Header().Get(echo.HeaderContentType)
			if resContentType == echo.MIMEApplicationJSON || resContentType == echo.MIMEApplicationJSONCharsetUTF8 {
				json.Unmarshal(resBodyBytes, &resBody)
			}
			logger.BodyDumpLog(c, resBody)
		},
	})
}

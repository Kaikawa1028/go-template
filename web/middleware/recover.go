package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/Kaikawa1028/go-template/app/logger"
	"net/http"
	"runtime"
)

// Recover PanicをRecoverするミドルウェアです
//
// 以下のファイルをベースにカスタマイズしています
// https://github.com/labstack/echo/blob/v4.5.0/middleware/recover.go
type Recover struct {
}

func NewRecover() *Recover {
	return &Recover{}
}

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func (m Recover) Recover() echo.MiddlewareFunc {
	return m.RecoverWithConfig(middleware.DefaultRecoverConfig)
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func (m Recover) RecoverWithConfig(config middleware.RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = middleware.DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					m.WritePanicLog(c, config, r)

					shouldResponse := !c.Response().Committed
					if shouldResponse {
						err := c.NoContent(http.StatusInternalServerError)
						if err != nil {
							logger.Error(c, err, nil)
						}
					}
				}
			}()
			return next(c)
		}
	}
}

func (m Recover) WritePanicLog(c echo.Context, config middleware.RecoverConfig, r interface{}) {
	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}
	stack := make([]byte, config.StackSize)
	length := runtime.Stack(stack, !config.DisableStackAll)
	if !config.DisablePrintStack {
		msg := fmt.Sprintf("%v %s\n", err, stack[:length])
		logger.Panic(c, msg, nil)
	}
}

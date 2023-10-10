package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/logger"
	"io/ioutil"
	"time"
)

// AccessLog RequestログとResponseログを出力するミドルウェア
type AccessLog struct {
}

func NewAccessLog() *AccessLog {
	return &AccessLog{}
}

func (m AccessLog) AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqBody := getRequestBody(c)
		logger.RequestLog(c, reqBody)

		var err error
		latency, latencyHuman := measureLatency(func() {
			err = next(c)
			if err != nil {
				c.Error(err)
			}
		})

		logger.ResponseLog(c, c.Response().Status, latency, latencyHuman)
		return nil
	}
}

func getRequestBody(c echo.Context) map[string]interface{} {
	var reqBody map[string]interface{}
	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if contentType == echo.MIMEApplicationJSON {
		if c.Request().Body != nil {
			reqBodyBytes, err := ioutil.ReadAll(c.Request().Body)
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes)) // Reset  see detail: https://stackoverflow.com/a/47295689
			if err == nil {
				json.Unmarshal(reqBodyBytes, &reqBody)
			}
		}
	}
	return reqBody
}

func measureLatency(proc func()) (latency time.Duration, latencyHuman string) {
	start := time.Now()
	proc()
	stop := time.Now()

	latency = stop.Sub(start)
	latencyHuman = latency.String()
	return
}

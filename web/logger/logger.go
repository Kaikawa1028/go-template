package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/errors"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

type (
	StackInfo struct {
		file     string
		line     int
		funcName string
	}
)

// ロガーの初期設定を行います
func SetupLogger() error {
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyMsg: "message",
		},
	})
	log.SetOutput(os.Stdout)

	level, err := getLogLevel()
	if err != nil {
		return errors.Wrap(err)
	}
	log.SetLevel(level)

	return nil
}

func Info(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		Info(msg)
}

func SimpleInfoF(format string, args ...interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, nil)).
		Infof(format, args...)
}

func Warn(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		Warn(msg)
}

func WarnWithError(c echo.Context, e error, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		WithField("panic", false).
		Warnf("%+v", e)
}

func Error(c echo.Context, e error, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		WithField("panic", false).
		Errorf("%+v", e)
}

func Debug(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		Debug(msg)
}

func Fatal(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		Fatal(msg)
}

func SimpleFatal(e error, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, params)).
		Fatalf("%+v", e)
}

func Panic(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		WithField("panic", true).
		Error(msg)
}

// Requestログを出力します
// AccessLogミドルウェア以外では使用しないでください
func RequestLog(c echo.Context, reqBody map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithFields(makeHttpFields(c)).
		WithField("input", reqBody).
		Info("[Request]")
}

// Responseログを出力します
// AccessLogミドルウェア以外では使用しないでください
func ResponseLog(c echo.Context, status int, latency time.Duration, latencyHuman string) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithFields(makeHttpFields(c)).
		WithField("latency", latency).
		WithField("latency_human", latencyHuman).
		WithField("status", status).
		Info("[Response]")
}

// Responseのボディをログに出力します
// BodyDumpLogミドルウェア以外では使用しないでください
func BodyDumpLog(c echo.Context, resBody map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	log.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithField("header", makeHeaderField(c)).
		WithField("output", resBody).
		Debug("[Body Dump]")
}

// ログレベルを取得します
func getLogLevel() (log.Level, error) {
	levelStr := os.Getenv("LOG_LEVEL")

	if levelStr == "" {
		return log.InfoLevel, nil
	}

	level, err := log.ParseLevel(levelStr)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return level, nil
}

// 各ログ共通のフィールドを組み立てます
func makeCommonFields(stackInfo *StackInfo, params map[string]interface{}) map[string]interface{} {
	var function *string
	var file *string
	var line *int
	if stackInfo != nil {
		function = &stackInfo.funcName
		file = &stackInfo.file
		line = &stackInfo.line
	}

	hostname, _ := os.Hostname()

	return map[string]interface{}{
		"params":   params,
		"function": function,
		"file":     file,
		"line":     line,
		"host":     hostname,
	}
}

// HTTPに関するフィールドを組み立てます
func makeHttpFields(c echo.Context) map[string]interface{} {
	req := c.Request()

	var userId *int
	if user, ok := c.Get("user").(*model.User); ok {
		userId = &user.ID
	}

	return map[string]interface{}{
		"uri":         req.RequestURI,
		"ip":          c.RealIP(),
		"http_method": req.Method,
		"server":      getLocalIP(),
		"referrer":    req.Referer(),
		"environment": os.Getenv("APP_ENV"),
		"header":      makeHeaderField(c),
		"userId":      userId,
	}
}

// echo.ContextからHTTPのヘッダーを取得します
// 併せて、一部のヘッダーの除外も行います
func makeHeaderField(c echo.Context) map[string]interface{} {
	excludeHeaders := []string{
		"Authorization",
	}
	return filterHeaders(c.Request().Header, excludeHeaders)
}

func makeStackInfo(pc uintptr, file string, line int, ok bool) *StackInfo {
	if !ok {
		return nil
	}

	funcName := runtime.FuncForPC(pc).Name()
	return &StackInfo{
		file:     file,
		line:     line,
		funcName: funcName,
	}
}

// 自サーバーのIPアドレスを取得します
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// HTTPのヘッダーから特定のヘッダーを除外して返します
func filterHeaders(headers http.Header, excludeHeaders []string) map[string]interface{} {
	if headers == nil {
		return nil
	}

	filteredHeaders := make(map[string]interface{})
	for k, v := range headers {
		if !includes(k, excludeHeaders) {
			filteredHeaders[k] = v
		}
	}
	return filteredHeaders
}

// 文字列の配列に特定の文字列が含まれているか調べます
func includes(target string, arr []string) bool {
	for _, v := range arr {
		if target == v {
			return true
		}
	}
	return false
}

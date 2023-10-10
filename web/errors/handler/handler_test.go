package handler_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/errors/handler"
	"github.com/Kaikawa1028/go-template/app/errors/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_標準エラーの場合500を返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	err := fmt.Errorf("standard error")
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, 0, rec.Body.Len())
	assertLogLevel(t, "error")
}

func TestHandler_ラップされた標準エラーの場合500を返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	err := errors.Wrap(fmt.Errorf("standard error"))
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, 0, rec.Body.Len())
	assertLogLevel(t, "error")
}

func TestHandler_権限エラーの場合403を返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	err := errors.Wrap(types.NewResourceNotPermittedError())
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, 0, rec.Body.Len())
	assertLogLevel(t, "warning")
}

func TestHandler_リソースが存在しないエラーの場合404を返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	err := errors.Wrap(types.NewResourceNotFoundError())
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, 0, rec.Body.Len())
	assertLogLevel(t, "warning")
}

func TestHandler_更新対象が存在しないエラーの場合404を返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	err := errors.Wrap(types.NewUpdateNoAffectedError())
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, 0, rec.Body.Len())
	assertLogLevel(t, "warning")
}

func TestHandler_Bindエラーの場合415を返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	err := errors.Wrap(echo.NewBindingError("", []string{""}, "", fmt.Errorf("")))
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	assert.Equal(t, 0, rec.Body.Len())
	assertLogLevel(t, "warning")
}

func TestHandler_Validationエラーの場合400を返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	err := errors.Wrap(validator.ValidationErrors{})
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	expectBody := `
{
	"errors": null,
	"method": "GET",
	"url": "/"
}`
	assert.JSONEq(t, expectBody, rec.Body.String())
	assertLogLevel(t, "warning")
}

func TestHandler_すでにレスポンスを返却済みのタイミングでエラーが発生した場合ログだけ記録されること(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 既にレスポンスを返却済み
	if e := c.NoContent(http.StatusOK); e != nil {
		t.Fatal(e)
	}

	err := fmt.Errorf("standard error")
	h := handler.NewErrorHandler()
	h.Handler(err, c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 0, rec.Body.Len())
	assertLogLevel(t, "error")
}

func TestHandler_HTTPErrorの場合エラーの構造体に保持しているステータスコードとメッセージを返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	err := errors.Wrap(echo.NewHTTPError(http.StatusNotFound, "dummy"))
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	expectBody := `{"message": "dummy"}`
	assert.JSONEq(t, expectBody, rec.Body.String())
	assertLogLevel(t, "warning")
}

func TestHandler_HTTPErrorの場合エラーの構造体に保持しているステータスコードとボディを返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	err := errors.Wrap(echo.NewHTTPError(http.StatusNotFound, map[string]string{"aaa": "bbb"}))
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	expectBody := `{"aaa": "bbb"}`
	assert.JSONEq(t, expectBody, rec.Body.String())
	assertLogLevel(t, "warning")
}

func TestHandler_HTTPErrorが更に別のHTTPErrorを保持する場合子のHTTPErrorに保持しているステータスコードとボディを返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	httpErr := echo.NewHTTPError(http.StatusForbidden, map[string]string{"ccc": "ddd"})
	httpErr.Internal = echo.NewHTTPError(http.StatusNotFound, map[string]string{"aaa": "bbb"})
	err := errors.Wrap(httpErr)
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	expectBody := `{"aaa": "bbb"}`
	assert.JSONEq(t, expectBody, rec.Body.String())
	assertLogLevel(t, "warning")
}

func TestHandler_HTTPErrorがHTTPError以外のエラーを保持する場合HTTPErrorに保持しているステータスコードとボディを返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	httpErr := echo.NewHTTPError(http.StatusForbidden, map[string]string{"ccc": "ddd"})
	httpErr.Internal = &json.UnmarshalTypeError{}
	err := errors.Wrap(httpErr)
	rec := callErrorHandler(err)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	expectBody := `{"ccc": "ddd"}`
	assert.JSONEq(t, expectBody, rec.Body.String())
	assertLogLevel(t, "warning")
}

func TestHandler_汎用業務エラーの場合422を返すこと(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)

	appErr := errors.Wrap(types.NewBasicBusinessError("テストメッセージ", map[string]interface{}{
		"testKey1": "testValue1",
	}))
	rec := callErrorHandler(appErr)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// レスポンス確認
	expectBody := `{"message": "テストメッセージ"}`
	assert.JSONEq(t, expectBody, rec.Body.String())

	// ログ確認
	assert.Len(t, loggerHook.Entries, 1)

	logStr, err := loggerHook.LastEntry().String()
	if err != nil {
		t.Fatal(err)
	}

	var log map[string]interface{}
	err = json.Unmarshal([]byte(logStr), &log)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "warning", log["level"])
	assert.Equal(t, "testValue1", log["params"].(map[string]interface{})["testKey1"])
}

func callErrorHandler(err error) *httptest.ResponseRecorder {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewErrorHandler()
	h.Handler(err, c)

	return rec
}

func assertLogLevel(t *testing.T, logLevel string) {
	assert.Len(t, loggerHook.Entries, 1)

	logStr, err := loggerHook.LastEntry().String()
	if err != nil {
		t.Fatal(err)
	}

	var log map[string]interface{}
	err = json.Unmarshal([]byte(logStr), &log)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, logLevel, log["level"])
}

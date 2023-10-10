package errors_test

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/errors/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestError_最下層のエラーメッセージを取得できること(t *testing.T) {
	// 組み込みエラーの場合
	standardErr := fmt.Errorf("ベースエラーメッセージ")
	wrappedErr1 := errors.Wrapf(standardErr, "追加メッセージ")
	wrappedErr2 := errors.Wrap(wrappedErr1)
	msg := wrappedErr2.Error()
	assert.Equal(t, "ベースエラーメッセージ", msg)

	// AppErrorの場合
	appErr := errors.New("ベースAppエラーメッセージ")
	wrappedErr1 = errors.Wrapf(appErr, "追加メッセージ")
	wrappedErr2 = errors.Wrap(wrappedErr1)
	msg = wrappedErr2.Error()
	assert.Equal(t, "ベースAppエラーメッセージ", msg)

	// Wrapしていない場合
	msg = appErr.Error()
	assert.Equal(t, "ベースAppエラーメッセージ", msg)
}

func TestGetRootError_Wrapされたエラーの場合最下層のエラーを取得できること(t *testing.T) {
	// 組み込みエラーの場合
	standardError := fmt.Errorf("組み込みエラー")
	wrappedError1 := errors.Wrap(standardError)
	wrappedError2 := errors.Wrap(wrappedError1)
	rootError := errors.GetRootError(wrappedError2)
	assert.Equal(t, standardError, rootError)

	// UpdateNoAffectedErrorの場合
	specialError := types.NewUpdateNoAffectedError()
	wrappedError1 = errors.Wrap(specialError)
	wrappedError2 = errors.Wrap(wrappedError1)
	rootError = errors.GetRootError(wrappedError2)
	assert.Equal(t, specialError, rootError)

	// Wrapしていない場合（組み込みエラー）
	rootError = errors.GetRootError(standardError)
	assert.Equal(t, standardError, rootError)

	// Wrapしていない場合（AppError）
	appErr := errors.New("test")
	rootError = errors.GetRootError(appErr)
	assert.Equal(t, appErr, rootError)

	// 外部ライブラリのエラーの場合（echo.HttpErrorの場合）
	httpError := echo.NewHTTPError(http.StatusNotFound, "")
	httpError.Internal = echo.NewHTTPError(http.StatusInternalServerError, "")
	wrappedError1 = errors.Wrap(httpError)
	rootError = errors.GetRootError(wrappedError1)
	assert.Equal(t, httpError, rootError)
}

func TestWrap_エラーをwrapしスタックトレースを取得できること(t *testing.T) {
	err1 := fmt.Errorf("ベースエラーメッセージ")
	err2 := errors.Wrapf(err1, "追加メッセージ")
	err3 := errors.Wrap(err2)

	s := fmt.Sprintf("%+v", err3)
	assert.Regexp(t, "ベースエラーメッセージ", s)
	assert.Regexp(t, "追加メッセージ", s)
	assert.Regexp(t, "errors/errors_test.go:[0-9]+", s)
}

//func TestAs_組み込みエラーとAppエラーを判別できること(t *testing.T) {
//	var target *errors.AppError
//
//	// 組み込みエラーの場合
//	standardError := fmt.Errorf("組み込みエラー")
//	isAppError := errors.As(standardError, &target)
//	assert.Equal(t, false, isAppError)
//	assert.Nil(t, target)
//
//	// Appエラーの場合
//	appError := errors.New("Appエラー")
//	isAppError = errors.As(appError, &target)
//	assert.Equal(t, true, isAppError)
//	assert.NotNil(t, target)
//}

func TestAsAppError_組み込みエラーとAppエラーを判別できること(t *testing.T) {
	// 組み込みエラーの場合
	standardError := fmt.Errorf("組み込みエラー")
	err := errors.AsAppError(standardError)
	assert.Nil(t, err)

	// Appエラーの場合
	appError := errors.New("Appエラー")
	err = errors.AsAppError(appError)
	assert.NotNil(t, err)
}

func TestRootAs_組み込みエラーとUpdateNoAffectedErrorを判別できること(t *testing.T) {
	var target types.UpdateNoAffectedError

	// 組み込みエラーの場合
	standardError := fmt.Errorf("組み込みエラー")
	wrappedError := errors.Wrap(standardError)
	isSpecialError := errors.RootAs(wrappedError, &target)
	assert.Equal(t, false, isSpecialError)

	// UpdateNoAffectedErrorの場合
	specialError := types.NewUpdateNoAffectedError()
	wrappedError = errors.Wrap(specialError)
	isSpecialError = errors.RootAs(wrappedError, &target)
	assert.Equal(t, true, isSpecialError)
	assert.Equal(t, "更新対象が存在していません。", target.Error())
}

package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/errors/types"
	"github.com/Kaikawa1028/go-template/app/interfaces/handler"
	"github.com/Kaikawa1028/go-template/app/logger"
	"github.com/Kaikawa1028/go-template/app/validate"
	"net/http"
)

// ErrorHandler 共通エラーハンドラ
// handlerからerrorをreturnした場合、このエラーハンドラで処理されます。
type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h ErrorHandler) Handler(err error, c echo.Context) {
	shouldResponse := !c.Response().Committed
	rootErr := errors.GetRootError(err)

	if _, ok := rootErr.(*echo.HTTPError); ok {
		// HTTPエラー
		logger.WarnWithError(c, err, nil)
		if shouldResponse {
			err = responseHttpError(rootErr.(*echo.HTTPError), c)
		}
	} else if _, ok := rootErr.(*echo.BindingError); ok {
		// バインドエラー
		logger.WarnWithError(c, err, nil)
		if shouldResponse {
			err = c.NoContent(http.StatusUnsupportedMediaType)
		}
	} else if _, ok := rootErr.(validator.ValidationErrors); ok {
		// バリデーションエラー
		logger.WarnWithError(c, err, nil)
		if shouldResponse {
			body := validate.MakeBadRequestErrorResponse(c, rootErr.(validator.ValidationErrors))
			err = c.JSON(http.StatusBadRequest, body)
		}
	} else if _, ok := rootErr.(types.ResourceNotPermittedError); ok {
		// 権限エラー
		logger.WarnWithError(c, err, nil)
		if shouldResponse {
			err = c.NoContent(http.StatusForbidden)
		}
	} else if _, ok := rootErr.(types.ResourceNotFoundError); ok {
		// リソースが存在しないエラー
		logger.WarnWithError(c, err, nil)
		if shouldResponse {
			err = c.NoContent(http.StatusNotFound)
		}
	} else if _, ok := rootErr.(types.UpdateNoAffectedError); ok {
		// 更新対象が存在しないエラー
		logger.WarnWithError(c, err, nil)
		if shouldResponse {
			err = c.NoContent(http.StatusNotFound)
		}
	} else if _, ok := rootErr.(types.BasicBusinessError); ok {
		// 汎用業務エラー
		businessErr := rootErr.(types.BasicBusinessError)
		logger.WarnWithError(c, err, businessErr.Params)
		if shouldResponse {
			body := h.makeBasicBusinessErrorResponse(businessErr)
			err = c.JSON(http.StatusUnprocessableEntity, body)
		}
	} else {
		// 上記以外のエラー
		logger.Error(c, err, nil)
		if shouldResponse {
			err = c.NoContent(http.StatusInternalServerError)
		}
	}

	// レスポンス返却時にエラーが発生した場合
	if shouldResponse && err != nil {
		logger.Error(c, err, nil)
	}
}

func (h ErrorHandler) makeBasicBusinessErrorResponse(businessError types.BasicBusinessError) *handler.BusinessErrorResponse {
	return &handler.BusinessErrorResponse{
		Message: businessError.Message,
	}
}

// responseHttpError echo.HTTPErrorを元にレスポンスを返します
// echoのDefaultHTTPErrorHandlerの実装を参考にしています
//   https://github.com/labstack/echo/blob/v4.5.0/echo.go#L360
func responseHttpError(he *echo.HTTPError, c echo.Context) (err error) {
	if he.Internal != nil {
		if herr, ok := he.Internal.(*echo.HTTPError); ok {
			he = herr
		}
	}

	code := he.Code
	message := he.Message
	if m, ok := he.Message.(string); ok {
		message = echo.Map{"message": m}
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(he.Code)
	} else {
		err = c.JSON(code, message)
	}
	return err
}

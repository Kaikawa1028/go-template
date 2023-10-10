package types

import "github.com/Kaikawa1028/go-template/app/errors"

//
// 業務エラーの構造体を定義しています
//

// BasicBusinessError 汎用的な業務エラー（エラーメッセージを返すだけの場合に使用します）
type BasicBusinessError struct {
	// Message エラーメッセージ（エンドユーザに表示されるほか、ログに出力されます）
	Message string

	// Params 任意のパラメータを格納します。ログに出力されます。
	Params map[string]interface{}
}

func NewBasicBusinessError(message string, params map[string]interface{}) error {
	return BasicBusinessError{
		Message: message,
		Params:  params,
	}
}

func (e BasicBusinessError) Error() string {
	return e.Message
}

func AsBasicBusinessError(err error) *BasicBusinessError {
	var target BasicBusinessError
	if matched := errors.RootAs(err, &target); matched {
		return &target
	} else {
		return nil
	}
}

// CantBase64DecodeAuthTokenError 認証トークンのBase64デコードに失敗
type CantBase64DecodeAuthTokenError struct {
}

func NewCantBase64DecodeAuthTokenError() error {
	return CantBase64DecodeAuthTokenError{}
}

func (e CantBase64DecodeAuthTokenError) Error() string {
	return "認証トークンのBase64デコードに失敗しました"
}

func AsCantBase64DecodeAuthTokenError(err error) *CantBase64DecodeAuthTokenError {
	var target CantBase64DecodeAuthTokenError
	if matched := errors.RootAs(err, &target); matched {
		return &target
	} else {
		return nil
	}
}

// NotFoundUserMatchedAuthTokenError 認証トークンに一致するユーザが見つからない
type NotFoundUserMatchedAuthTokenError struct {
}

func NewNotFoundUserMatchedAuthTokenError() error {
	return NotFoundUserMatchedAuthTokenError{}
}

func (e NotFoundUserMatchedAuthTokenError) Error() string {
	return "認証トークンに一致するユーザが見つかりませんでした"
}

func AsNotFoundUserMatchedAuthTokenError(err error) *NotFoundUserMatchedAuthTokenError {
	var target NotFoundUserMatchedAuthTokenError
	if matched := errors.RootAs(err, &target); matched {
		return &target
	} else {
		return nil
	}
}

// CmsItemCsvImportValidationError cms_itemsのCSVインポートのバリデーションエラー
type CmsItemCsvImportValidationError struct {
	Message string
}

func NewCmsItemValueCsvImportValidationError(message string) error {
	return CmsItemCsvImportValidationError{
		message,
	}
}

func (e CmsItemCsvImportValidationError) Error() string {
	return e.Message
}

func AsCmsItemValueCsvImportValidationError(err error) *CmsItemCsvImportValidationError {
	var target CmsItemCsvImportValidationError
	if matched := errors.RootAs(err, &target); matched {
		return &target
	} else {
		return nil
	}
}

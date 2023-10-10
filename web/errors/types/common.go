package types

import (
	"github.com/Kaikawa1028/go-template/app/errors"
)

//
// 汎用的に使用するエラーの構造体を定義しています
//

// ResourceNotFoundError リソースが存在しない場合のエラー
type ResourceNotFoundError struct{}

func NewResourceNotFoundError() error {
	return ResourceNotFoundError{}
}

func (e ResourceNotFoundError) Error() string {
	return "リソースが存在していません。"
}

func AsResourceNotFoundError(err error) *ResourceNotFoundError {
	var target ResourceNotFoundError
	if matched := errors.RootAs(err, &target); matched {
		return &target
	} else {
		return nil
	}
}

type ResourceNotPermittedError struct{}

func NewResourceNotPermittedError() error {
	return ResourceNotPermittedError{}
}

func (e ResourceNotPermittedError) Error() string {
	return "リソースへの権限がありません。"
}

func AsResourceNotPermittedError(err error) *ResourceNotPermittedError {
	var target ResourceNotPermittedError
	if matched := errors.RootAs(err, &target); matched {
		return &target
	} else {
		return nil
	}
}

// UpdateNoAffectedError DB上に更新対象が存在しない場合のエラー
type UpdateNoAffectedError struct{}

func NewUpdateNoAffectedError() error {
	return UpdateNoAffectedError{}
}

func (e UpdateNoAffectedError) Error() string {
	return "更新対象が存在していません。"
}

func AsUpdateNoAffectedError(err error) *UpdateNoAffectedError {
	var target UpdateNoAffectedError
	if matched := errors.RootAs(err, &target); matched {
		return &target
	} else {
		return nil
	}
}

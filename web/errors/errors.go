package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

// AppError エラー構造体
type AppError struct {
	// 標準のエラー仕様を満たす変数
	next    error
	message string
	frame   xerrors.Frame
}

// Error 最下層のエラーメッセージを取得する
func (e *AppError) Error() string {
	rootErr := GetRootError(e)
	if rootAppErr := AsAppError(rootErr); rootAppErr != nil {
		return rootAppErr.message
	} else {
		return rootErr.Error()
	}
}

// GetRootError 最下層のエラーを取り出す
func GetRootError(err error) error {
	appErr := AsAppError(err)
	if appErr != nil {
		if appErr.next != nil {
			return GetRootError(appErr.next)
		} else {
			return appErr
		}
	} else {
		return err
	}
}

// Format 書式指定に従ってエラーメッセージを組み立てます
func (e *AppError) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }

// FormatError 書式指定に従ってエラーメッセージを組み立てます
func (e *AppError) FormatError(p xerrors.Printer) error {
	var message string
	if e.message != "" {
		message += fmt.Sprintf("%s", e.message)
	}

	p.Print(message)
	e.frame.Format(p)
	return e.next
}

func create(msg string) *AppError {
	var e AppError
	e.message = msg
	e.frame = xerrors.Caller(2)

	return &e
}

// New エラーを作成する
func New(msg string) *AppError {
	return create(msg)
}

// Errorf エラーを作成する
func Errorf(format string, args ...interface{}) *AppError {
	return create(fmt.Sprintf(format, args...))
}

// Wrap エラーを内包した新しいエラーを作成します
// エラーを内包すると同時に、スタックトレースを出力するのに必要な情報（ファイル名や行番号）を収集します
func Wrap(err error, msg ...string) *AppError {
	if err == nil {
		return nil
	}

	var m string
	if len(msg) != 0 {
		m = msg[0]
	}
	e := create(m)
	e.next = err
	return e
}

// Unwrap 子のエラーを返します
func (e *AppError) Unwrap() error { return e.next }

// Wrapf エラーをラップして新しいエラーを作成する（メッセージ付き）
func Wrapf(err error, format string, args ...interface{}) *AppError {
	e := create(fmt.Sprintf(format, args...))
	e.next = err
	return e
}

// RootAs 最下層のエラーがtarget(第二引数)のタイプと一致するかチェックし
// 一致する場合はtargetにその値を設定してtrueを返す
func RootAs(err error, target interface{}) bool {
	rootErr := GetRootError(err)
	return xerrors.As(rootErr, target)
}

// AsAppError エラーがAppErrorの場合は、その参照を返す
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var e *AppError
	if xerrors.As(err, &e) {
		return e
	}
	return nil
}

package a_err

import "fmt"

type CustomError struct {
	Status int
	Code   string
	Type   string
	Msg    string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Status: %d, Code: %s, Msg: %s Type: %s", e.Status, e.Code, e.Type, e.Msg)
}

func Throw(err *CustomError) {
	panic(err)
}

func NewServerErr(msg string) *CustomError {
	return &CustomError{
		Status: 500,
		Code:   "E0001",
		Type:   "サーバーエラー",
		Msg:    msg,
	}
}
func NewDBErr(msg string) *CustomError {
	return &CustomError{
		Status: 500,
		Code:   "E0002",
		Type:   "DBエラー",
		Msg:    msg,
	}
}
func NewRequestErr(msg string) *CustomError {
	return &CustomError{
		Status: 400,
		Code:   "W0001",
		Type:   "リクエストエラー",
		Msg:    msg,
	}
}
func NewNotFoundErr(msg string) *CustomError {
	return &CustomError{
		Status: 404,
		Code:   "W0004",
		Type:   "非存在エラー",
		Msg:    msg,
	}
}
func NewAuthErr(msg string) *CustomError {
	return &CustomError{
		Status: 401,
		Code:   "W0002",
		Type:   "認証エラー",
		Msg:    msg,
	}
}

var InternalServerErr = &CustomError{
	500,
	"E0001",
	"サーバーエラーが発生しました",
	"",
}
var BadRequestErr = &CustomError{
	400,
	"W0001",
	"リクエストエラーが発生しました",
	"",
}
var UnauthorizedErr = &CustomError{
	401,
	"W0002",
	"認証エラーが発生しました",
	"",
}
var ForbiddenErr = &CustomError{
	403,
	"W0003",
	"権限がありません",
	"",
}
var NotFoundErr = &CustomError{
	404,
	"W0004",
	"該当の項目が見つかりませんでした",
	"",
}

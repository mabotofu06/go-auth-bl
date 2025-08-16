package ctm_err

// go-auth-blで使用するカスタムエラー定義
type CustomError struct {
	Status int
	Code   string
	Type   string
	Msg    string
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
func NewAuthErr(msg string) *CustomError {
	return &CustomError{
		Status: 401,
		Code:   "W0002",
		Type:   "認証エラー",
		Msg:    msg,
	}
}
func NewPermissionErr(msg string) *CustomError {
	return &CustomError{
		Status: 403,
		Code:   "W0003",
		Type:   "権限エラー",
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

var (
	InternalServerErr   = NewServerErr("サーバーエラーが発生しました")
	UnexpectedServerErr = NewServerErr("予期せぬエラーが発生しました")

	UnexpectedDBErr = NewDBErr("予期せぬDBエラーが発生しました")

	BadRequestErr = NewRequestErr("リクエストエラーが発生しました")
	ParameterErr  = NewRequestErr("パラメータエラーが発生しました")

	UnauthorizedErr = NewAuthErr("認証エラーが発生しました")

	ForbiddenErr = NewPermissionErr("権限エラーが発生しました")

	NotFoundErr = NewNotFoundErr("非存在エラーが発生しました")
)

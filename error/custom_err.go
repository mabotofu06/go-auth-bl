package a_err

import "fmt"

type CustomError struct {
	Status int
	Code   string
	Msg    string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Status: %d, Code: %s, Msg: %s", e.Status, e.Code, e.Msg)
}

var InternalServerErr = &CustomError{
	500,
	"E0001",
	"Internal Server Error",
}
var BadRequestErr = &CustomError{
	400,
	"W0001",
	"Bad Request",
}
var UnauthorizedErr = &CustomError{
	401,
	"W0002",
	"Unauthorized",
}
var ForbiddenErr = &CustomError{
	403,
	"W0003",
	"Forbidden",
}
var NotFoundErr = &CustomError{
	404,
	"W0004",
	"Not Found",
}

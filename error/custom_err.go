package a_err

import "fmt"

type CustomError struct {
	code    int
	message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
}

var NotFoundErr = &CustomError{404, "Not Found"}
var InternalServerErr = &CustomError{500, "Internal Server Error"}
var BadRequestErr = &CustomError{400, "Bad Request"}
var UnauthorizedErr = &CustomError{401, "Unauthorized"}
var ForbiddenErr = &CustomError{403, "Forbidden"}

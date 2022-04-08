package errors

import (
	"fmt"
	"net/http"
)

func WithCode(err error, code int, message string) *Error {
	return &Error{
		code:    code,
		message: message,
		cause:   err,
	}
}

func AsCode(err error) *Error {
	e := new(Error)
	if As(err, &e) {
		return e
	}
	return &Error{
		code: http.StatusInternalServerError,
		message: "unknown",
		cause: err,
	}
}

func Code(err error) int {
	 e := new(Error)
	if As(err, &e) {
		return e.code
	}
	return 0
}

func Message(err error) string {
	e := new(Error)
	if As(err, &e) {
		return e.message
	}
	return ""
}

// 大类错误
type Error struct {
	code    int
	message string
	cause   error
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Error() string {
	return fmt.Sprintf("code = %d, message = %s, cause = %v", e.code, e.message, e.cause)
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) Equal(err error) bool {
	return e.Code() == Code(err)
}

package coderror

import (
	"errors"
	"fmt"
)

type Code interface {
	Int() int
	String() string
}

type errorWithCode struct {
	code Code
	err  error
}

func New(code Code, err error) error {
	return &errorWithCode{
		code: code,
		err:  err,
	}
}

func Newf(code Code, format string, arg ...interface{}) error {
	return New(code, fmt.Errorf(format, arg...))
}

func (e errorWithCode) Error() string {
	return e.err.Error()
}

func (e errorWithCode) Code() Code {
	return e.code
}

func Errorf(err error, format string, arg ...interface{}) error {
	if cerr := (errorWithCode{}); errors.As(err, &cerr) {
		cerr.err = fmt.Errorf(format, arg...)
		return cerr
	}

	return fmt.Errorf(format, arg...)
}

func As(err error) (*errorWithCode, bool) {
	if cerr := (errorWithCode{}); errors.As(err, &cerr) {
		return &cerr, true
	}

	return nil, false
}

package errors

import (
	"github.com/pkg/errors"
)

type UserError struct {
	ErrorCode    int
	ErrorMessage string
}
type Error struct {
	Error     error     `json:"err"`
	UserError UserError `json:"usererror"`
}

func Wrap(e error, s string) error {
	return errors.Wrap(e, s)
}
func Wrapf(e error, s string, args ...interface{}) error {
	return errors.Wrapf(e, s, args)
}

func NewError(s string) error {
	return errors.New(s)
}

func New() *Error {
	var e Error
	return &e
}

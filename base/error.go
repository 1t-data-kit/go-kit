package base

import (
	"fmt"
	"github.com/pkg/errors"
)

type Error struct {
	code  interface{}
	human string
	err   error
	wrap  bool
}

func NewError(code interface{}, human string) *Error {
	return &Error{
		code:  code,
		human: human,
		wrap:  true,
	}
}

func (e *Error) WithUnWrap() *Error {
	_err := e.Clone()
	_err.wrap = false
	return _err
}

func (e *Error) WithError(err error) *Error {
	_err := e.Clone()
	_err.err = err
	return _err
}

func (e *Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

func (e *Error) Code() interface{} {
	return e.code
}

func (e *Error) Human() string {
	return e.human
}

func (e *Error) Clone() *Error {
	_e := NewError(e.code, e.human)
	_e.err = e.err
	_e.wrap = e.wrap
	return _e
}

func (e *Error) Wrap(err error, message string) *Error {
	if !ErrorCanWrap(err) {
		return err.(*Error)
	}
	_err := e.Clone()
	_err.err = errors.Wrap(err, message)
	return _err
}

func (e *Error) Wrapf(err error, format string, args ...interface{}) *Error {
	if !ErrorCanWrap(err) {
		return err.(*Error)
	}
	_err := e.Clone()
	_err.err = errors.Wrapf(err, format, args...)
	return _err
}

func (e *Error) ReHuman(human string) *Error {
	_err := e.Clone()
	_err.human = human
	return _err
}

func (e *Error) ReHumanf(format string, args ...interface{}) *Error {
	_err := e.Clone()
	_err.human = fmt.Sprintf(format, args...)
	return _err
}

func IsError(err error) bool {
	if _, ok := err.(*Error); ok {
		return true
	}
	return false
}

func ErrorCanWrap(err error) bool {
	if _e, ok := err.(*Error); ok {
		return _e.wrap
	}
	return true
}

func ErrorCode(err error) (interface{}, bool) {
	if e, ok := err.(*Error); ok {
		return e.code, true
	}
	return nil, false
}

func ErrorMessage(err error) string {
	if e, ok := err.(*Error); ok {
		return e.Error()
	}
	return err.Error()
}

func ErrorHuman(err error) string {
	if e, ok := err.(*Error); ok {
		return e.human
	}
	return ""
}

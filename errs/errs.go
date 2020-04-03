// Package errs exports New, Errorf, NoWrap and CloseIgnoreErr.
package errs

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// CloseIgnoreErr closes c ignoring the error. It is intended to be used
// when an io.Closer is closed in a defer statement to avoid a lint
// error about ignoring the return value. The name of the function makes
// it clear that the error is being ignored and is less cumbersome than
// adding a `//nolint:errcheck` or using an anonymous func in defer:
// `defer func() { _ = c.Close() }()`
func CloseIgnoreErr(c io.Closer) {
	_ = c.Close()
}

// New wraps multiple errors and formats them as
//	err1: err2: err3 ...
func New(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	format := strings.Repeat("%v: ", len(errs)-1) + "%v"
	args := make([]interface{}, len(errs))
	for i := range errs {
		args[i] = errs[i]
	}
	return Errorf(format, args...)
}

// Errorf works in analogy to fmt.Errorf only that it wraps all errors
// except for the ones marked as NoWrap.
func Errorf(format string, args ...interface{}) error {
	rawArgs := make([]interface{}, len(args))
	var errs []error
	for i, arg := range args {
		if n, ok := arg.(noWrap); ok {
			arg = n.err
		} else if e, ok := arg.(error); ok {
			errs = append(errs, e)
		}
		rawArgs[i] = arg
	}
	return &multiErr{
		s:    fmt.Sprintf(format, rawArgs...),
		errs: errs,
	}
}

// NoWrap marks errors not to be wrapped for Errorf.
func NoWrap(err error) noWrap { //nolint:golint
	return noWrap{err: err}
}

type noWrap struct{ err error }

// multiErr contains wrapped errors and a formatted error message.
type multiErr struct {
	s    string
	errs []error
}

// Error returns multiErr's error message to implement error interface.
func (e *multiErr) Error() string {
	return e.s
}

// Is reports whether any error in mutiErr's slice of errors or any
// error within the error's chain matches target.
func (e *multiErr) Is(target error) bool {
	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// As finds the first error in multiErr's slice of errors or any error
// within the error's chain that matches target, and if so, sets target
// to that error value and returns true.
func (e *multiErr) As(target interface{}) bool {
	for _, err := range e.errs {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

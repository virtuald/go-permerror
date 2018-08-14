// Package permerror creates errors that have a `Temporary` function to be
// used with grpc's `FailOnNonTempDialError` option.
//
// Designed in the spirit of github.com/pkg/errors, the returned errors all
// implement the non-exported causer interface.
package permerror

import (
	"github.com/pkg/errors"
)

// MakePermanent forces an error to be permanent
func MakePermanent(cause error) error {
	return &madePermanent{cause: cause}
}

// New returns an error message and marks it as permanent
func New(msg string) error {
	return &permError{msg: msg}
}

// Wrap wraps an error and marks it as permanent unless the
// underlying error says otherwise
func Wrap(cause error) error {
	return &wrapError{cause: cause}
}

// WithMessage wraps an error and marks it as permanent unless the
// underlying error says otherwise
func WithMessage(cause error, msg string) error {
	return &permErrorWrapper{
		cause: cause,
		msg:   msg,
	}
}

type permErrorWrapper struct {
	cause error
	msg   string
}

func (pe *permErrorWrapper) Error() string { return pe.msg + ": " + pe.cause.Error() }
func (pe *permErrorWrapper) Cause() error  { return pe.cause }

func (pe *permErrorWrapper) Temporary() bool {
	err := errors.Cause(pe.cause)
	switch err := err.(type) {
	case interface {
		Temporary() bool
	}:
		return err.Temporary()
	default:
		// default to permanent if not specified by the cause
		return false
	}
}

type permError struct {
	msg string
}

func (pe *permError) Error() string { return pe.msg }
func (*permError) Temporary() bool  { return false }

type madePermanent struct {
	cause error
}

func (mp *madePermanent) Error() string { return mp.cause.Error() }
func (mp *madePermanent) Cause() error  { return mp.cause }
func (*madePermanent) Temporary() bool  { return false }

type wrapError struct {
	cause error
}

func (we *wrapError) Error() string { return we.cause.Error() }
func (we *wrapError) Cause() error  { return we.cause }

func (we *wrapError) Temporary() bool {
	err := errors.Cause(we.cause)
	switch err := err.(type) {
	case interface {
		Temporary() bool
	}:
		return err.Temporary()
	default:
		// default to permanent if not specified by the cause
		return false
	}
}

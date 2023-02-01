package errors

import (
	"fmt"
	"github.com/eadydb/hubble/proto/v1"
)

type Error interface {
	Error() string
	StatusCode() proto.StatusCode
	Suggestions() []*proto.Suggestion
	Unwrap() error
}

type ErrDef struct {
	err error
	ae  *proto.ActionableErr
}

var _ error = (*ErrDef)(nil)

func (e *ErrDef) Error() string {
	return fmt.Sprintf("xx")
}

func (e *ErrDef) Unwrap() error {
	return e.err
}

func (e *ErrDef) StatusCode() proto.StatusCode {
	return e.ae.ErrCode
}

func (e *ErrDef) Suggestions() []*proto.Suggestion {
	return e.ae.Suggestions
}

// NewError creates an actionable error message preserving the actual error.
func NewError(err error, ae *proto.ActionableErr) *ErrDef {
	return &ErrDef{
		err: err,
		ae:  ae,
	}
}

// NewErrorWithStatusCode creates an actionable error message.
func NewErrorWithStatusCode(ae *proto.ActionableErr) *ErrDef {
	return &ErrDef{
		ae: ae,
	}
}

func IsHubbleErr(err error) bool {
	if _, ok := err.(Error); ok {
		return true
	}
	return false
}

package uerror

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// NoType - basic err type
	NoType = ErrorType(iota)
	// NotFound - entity not found
	NotFound
	// ServerError - something went wrong
	ServerError
	// NotAuthorized - user has no credentials
	NotAuthorized
	// BadRequest - validation or type error
	BadRequest
	// Forbidden - credentials present, but not satisfy access control
	Forbidden
)

// stackTracer - interface for checking stack trace availability
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// ErrorType - error types classifier
type ErrorType uint
type errorContext map[string]string

type customError struct {
	errorType     ErrorType
	originalError error
	context       errorContext
}

func (error customError) Error() string {
	return error.originalError.Error()
}

// New - custom error constructor
func (t ErrorType) New(message string) error {
	return customError{
		errorType:     t,
		originalError: errors.New(message),
	}
}

// Newf - custom error formatted constructor
func (t ErrorType) Newf(message string, args ...interface{}) error {
	err := fmt.Errorf(message, args...)
	return customError{
		errorType:     t,
		originalError: err,
	}
}

// Wrapf - custom error formatted wrapper
func (t ErrorType) Wrapf(err error, message string, args ...interface{}) error {
	nErr := errors.Wrapf(err, message, args...)
	return customError{
		errorType:     t,
		originalError: nErr,
	}
}

// Wrap - custom error wrapper
func (t ErrorType) Wrap(err error, message string) error {
	return t.Wrapf(err, message)
}

// New - returns new NoType custom error
func New(msg string) error {
	return customError{
		errorType:     NoType,
		originalError: errors.New(msg),
	}
}

// Newf - returns new NoType formatted custom error
func Newf(msg string, args ...interface{}) error {
	return customError{}
}

// Wrapf - wraps an error (simple or custom) into custom error with stack and context formatted
func Wrapf(err error, message string, args ...interface{}) error {
	wrappedErr := errors.Wrapf(err, message, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			originalError: wrappedErr,
			context:       customErr.context,
		}
	}
	return customError{
		errorType:     NoType,
		originalError: wrappedErr,
	}
}

// Wrap - wraps an error (simple or custom) into custom error with stack and context
func Wrap(err error, message string) error {
	return Wrapf(err, message)
}

// Cause - returns an original error
func Cause(err error) error {
	return errors.Cause(err)
}

// AddContext - adds a key-value pair to error's context
func AddContext(err error, field, message string) error {
	if customErr, ok := err.(customError); ok {
		if customErr.context == nil {
			customErr.context = errorContext{}
		}
		customErr.context[field] = message
		return customErr
	}
	return customError{
		errorType:     NoType,
		originalError: err,
		context:       errorContext{field: message},
	}
}

// GetContext - returns an attached error context if possible
func GetContext(err error) map[string]string {
	if customErr, ok := err.(customError); ok {
		return map[string]string(customErr.context)
	}
	return nil
}

// GetStack - returns stack trace as string as possible, with explicit depth
func GetStack(err error, depth int) string {
	if customErr, ok := err.(customError); ok {
		err = customErr.originalError
	}
	if err, ok := err.(stackTracer); ok {
		return fmt.Sprintf("%+v", err.StackTrace()[0:depth])
	}
	return ""
}

// GetType - returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}
	return NoType
}

package types

import (
	"context"
	"errors"
	"log"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Errors
var (
	ErrWrongPassword      = errors.New("wrong password")
	ErrWrongEmail         = errors.New("wrong email")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNotFound           = errors.New("not found")
	ErrNoInput            = errors.New("no input")
	ErrLimitInput         = errors.New("name should be more than 5 char")
	ErrNameAlreadyExist   = errors.New("name already exits")
	ErrClerkValidationErr = errors.New("clerk validation error")
)

var (
	ErrTypesHandlerError = "handler-error"
	ErrTypesRepoError    = "repo-error"
	ErrTypesServiceError = "service-error"
	ErrTypesClerkError   = "clerk-error"
)

// Error represents customized error object
type Error struct {
	Path     string
	Message  string
	Error    error
	Type     string
	IsIgnore bool
	FuncName string
	File     string
	Line     int
}

// ✅ Centralized logging and tracing
func (e *Error) Log(ctx context.Context, tracer trace.Tracer) {
	if e == nil || e.Error == nil {
		return
	}

	funcName := e.FuncName
	file := e.File
	line := e.Line

	log.Printf("[ERROR] %s:%d %s | %s: %v\n", file, line, funcName, e.Message, e.Error)

	_, span := tracer.Start(ctx, "error:"+funcName)
	defer span.End()

	span.RecordError(e.Error)
	span.SetAttributes(
		attribute.String("error.path", e.Path),
		attribute.String("error.type", e.Type),
		attribute.String("error.message", e.Message),
		attribute.String("error.func", funcName),
		attribute.String("error.file", file),
		attribute.Int("error.line", line),
	)
}

func (e *Error) LogAndReturn(ctx context.Context, tracer trace.Tracer) *Error {
	e.Log(ctx, tracer)
	return e
}

func NewError(err error) *Error {
	pc, file, line, ok := runtime.Caller(1) // caller of NewError
	funcName := "unknown"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = fn.Name()
		}
	}

	return &Error{
		// Path:     path,
		Message: err.Error(),
		// Type:     errType,
		Error:    err,
		FuncName: funcName,
		File:     file,
		Line:     line,
	}
}

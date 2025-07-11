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
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNotFound           = errors.New("not found")
	ErrNoInput            = errors.New("no input")
	ErrLimitInput         = errors.New("name should be more than 5 char")
	ErrNameAlreadyExist   = errors.New("name already exits")
	ErrClerkValidationErr = errors.New("clerk validation error")
)

// Error represents customized error object
type Error struct {
	Path     string
	Message  string
	Error    error
	Type     string
	IsIgnore bool
}

// ✅ Centralized logging and tracing
func (e *Error) Log(ctx context.Context, tracer trace.Tracer) {
	if e == nil || e.Error == nil {
		return
	}

	pc, file, line, ok := runtime.Caller(1)
	funcName := "unknown"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = fn.Name()
		}
	}

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

package response

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/riskibarqy/bq-account-service/external/logger"
	"github.com/riskibarqy/bq-account-service/internal/types"
	"gopkg.in/go-playground/validator.v9"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Fields  []*FieldError `json:"fields,omitempty"`
}

func MakeFieldError(field, tag string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: fmt.Sprintf("Validation failed on tag '%s'", tag),
	}
}

// Error writes error http response and logs via Uptrace + terminal
func Error(w http.ResponseWriter, message string, status int, err types.Error) {
	// Step 1: Log error via logger + uptrace
	err.Log(context.TODO(), logger.Tracer)

	// Step 2: Setup basic headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Step 3: Choose error code
	errorCode := "InternalServerError"
	switch status {
	case http.StatusUnauthorized:
		errorCode = "Unauthorized"
	case http.StatusNotFound:
		errorCode = "NotFound"
	case http.StatusBadRequest:
		errorCode = "BadRequest"
	case http.StatusUnprocessableEntity:
		errorCode = "ValidationError"
	}

	// Step 4: Handle validation errors
	errorFields := []*FieldError{}
	if ve, ok := err.Error.(validator.ValidationErrors); ok {
		message = "Validation failed"
		for _, fieldErr := range ve {
			errorFields = append(errorFields, MakeFieldError(fieldErr.Field(), fieldErr.ActualTag()))
		}
	}

	// Step 5: Encode response
	res := ErrorResponse{
		Code:    errorCode,
		Message: message,
		Fields:  errorFields,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("[response.Error] failed to encode JSON: %v", err)
	}

	// Step 6: Optional stack trace logging (e.g. from pkg/errors)
	if err.Error != nil {
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}
		if e, ok := err.Error.(stackTracer); ok {
			st := e.StackTrace()
			if len(st) > 0 {
				fmt.Printf("[STACKTRACE] %s: %+v\n", err.Message, st[0])
			}
		}
	}
}

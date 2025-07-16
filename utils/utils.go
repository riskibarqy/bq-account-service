package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"
	"unicode"

	"github.com/riskibarqy/bq-account-service/config"
	"github.com/riskibarqy/bq-account-service/external/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Now() int {
	return int(time.Now().Unix())
}

func EncodeHexMD5(params string) string {
	sumString := md5.Sum([]byte(params))
	return hex.EncodeToString(sumString[:])
}

func SplitName(fullName string) (firstName, lastName string) {
	parts := strings.Fields(fullName) // split by whitespace

	if len(parts) == 0 {
		return "", ""
	}

	firstName = parts[0]

	if len(parts) == 1 {
		lastName = ""
	} else {
		lastName = strings.Join(parts[1:], " ")
	}

	return firstName, lastName
}

// CreateUsernameFromEmail returns a sanitized username generated from an email address.
// e.g., "john.doe@example.com" → "john_doe"
func CreateUsernameFromEmail(email string) string {
	// Extract the part before @
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return "invalid_username"
	}

	prefix := email[:atIndex]

	// Replace dots with underscores, remove invalid characters
	var builder strings.Builder
	for _, ch := range prefix {
		switch {
		case unicode.IsLetter(ch), unicode.IsDigit(ch):
			builder.WriteRune(ch)
		case ch == '.' || ch == '-' || ch == '_':
			builder.WriteRune('_')
			// ignore other characters
		}
	}

	username := builder.String()

	// Ensure it’s not empty
	if username == "" {
		return "user"
	}

	return username
}

func WithDBSpan(ctx context.Context, operation string, statement string, fn func(ctx context.Context) error) error {
	tracer := logger.Tracer
	ctx, span := tracer.Start(ctx, operation,
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.name", config.AppConfig.DBName),
			attribute.String("db.statement", statement),
		),
	)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

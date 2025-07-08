package logger

import (
	"context"
	"log"
	"runtime"
	"strings"

	"github.com/riskibarqy/bq-account-service/config"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// Call this once in main()
func Init() {

	// Configure Uptrace with DSN and optional tags
	uptrace.ConfigureOpentelemetry(
		// Make sure UPTRACE_DSN is set in env
		// e.g., export UPTRACE_DSN=https://<token>@api.uptrace.dev/project_id
		uptrace.WithDSN(config.AppConfig.UptraceDSN),
		uptrace.WithServiceName(config.AppConfig.AppName),
		uptrace.WithServiceVersion(config.AppConfig.AppVersion),
		uptrace.WithDeploymentEnvironment(config.AppConfig.AppMode),
	)

	tracer = otel.Tracer("logger")
}

// Call this at shutdown, e.g. defer logger.Shutdown(ctx)
func Shutdown(ctx context.Context) {
	_ = uptrace.Shutdown(ctx)
}

type LogLevel string

const (
	LevelDebug   LogLevel = "DEBUG"
	LevelInfo    LogLevel = "INFO"
	LevelWarn    LogLevel = "WARN"
	LevelError   LogLevel = "ERROR"
	LevelSuccess LogLevel = "SUCCESS"
	LevelFatal   LogLevel = "FATAL"
)

func Log(ctx context.Context, level LogLevel, message string, err error, fields ...attribute.KeyValue) {
	// Get caller info
	pc, file, line, ok := runtime.Caller(1)
	funcName := "unknown"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = fn.Name()
		}
	}
	shortFile := file
	if i := strings.LastIndex(file, "/"); i != -1 {
		shortFile = file[i+1:]
	}

	// Log to terminal
	log.Printf("[%s] %s:%d %s: %s %v\n", level, shortFile, line, funcName, message, err)

	// Only send to Uptrace for error/warn/fatal
	if level == LevelError || level == LevelWarn || level == LevelFatal {
		_, span := tracer.Start(ctx, string(level)+":"+funcName)
		defer span.End()

		if err != nil {
			span.RecordError(err)
		}

		allAttrs := []attribute.KeyValue{
			attribute.String("log.level", string(level)),
			attribute.String("log.message", message),
			attribute.String("file", file),
			attribute.String("func", funcName),
			attribute.Int("line", line),
		}
		allAttrs = append(allAttrs, fields...)

		span.SetAttributes(allAttrs...)
	}
}

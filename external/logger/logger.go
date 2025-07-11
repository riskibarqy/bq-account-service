package logger

import (
	"context"
	"log"

	"github.com/riskibarqy/bq-account-service/config"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

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

	Tracer = otel.Tracer("logger")
	log.Println("[Uptrace] Initialized")
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

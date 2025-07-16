package databases

import (
	"log"

	"github.com/riskibarqy/bq-account-service/config"
	sqlxtrace "github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	sqltrace "go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Init Connect to the database
func Init() {
	_, err := sqltrace.Register("postgres",
		sqltrace.AllowRoot(),
		sqltrace.TraceQueryWithoutArgs(),
		sqltrace.TraceRowsClose(),
		sqltrace.TraceRowsAffected(),
		sqltrace.WithDatabaseName(config.AppConfig.DBName), // Optional.
		sqltrace.WithSystem(semconv.DBSystemPostgreSQL),    // Optional.
	)
	if err != nil {
		log.Fatalf("Failed Register SQLTrace : %v", err)
	}
	// Open new database connection
	db, err := sqlxtrace.Open("postgres", config.AppConfig.DBConnectionString)
	if err != nil {
		log.Fatalf("Failed to reconnect to the database: %v", err)
	}

	// Assign new connection to AppConfig
	config.AppConfig.DatabaseClient = db

	log.Println("[Postgres] Successfully connected to the database")
}

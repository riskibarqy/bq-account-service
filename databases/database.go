package databases

import (
	"log"

	"github.com/XSAM/otelsql"
	"github.com/riskibarqy/bq-account-service/config"
	otelsqlx "github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	"go.opentelemetry.io/otel/attribute"
)

func Init() {
	// Register otelsql driver with additional span attributes
	_, err := otelsql.Register("postgres",
		otelsql.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.name", config.AppConfig.DBName),
		),
		otelsql.WithSQLCommenter(true), // Adds trace info as SQL comment for debugging
	)
	if err != nil {
		log.Fatalf("Failed to register otelsql driver: %v", err)
	}

	db, err := otelsqlx.Open("postgres", config.AppConfig.DBConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	config.AppConfig.DatabaseClient = db

	log.Println("[Postgres] Connected with OpenTelemetry instrumentation")
}

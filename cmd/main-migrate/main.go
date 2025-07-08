package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/riskibarqy/bq-account-service/config"
	"github.com/riskibarqy/bq-account-service/databases"
)

func main() {
	config.GetConfiguration()
	databases.NewMigrator(config.AppConfig.DBConnectionString).Run("up")
}

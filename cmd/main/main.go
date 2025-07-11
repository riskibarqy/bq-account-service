package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ancalabrese/reload"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/riskibarqy/bq-account-service/config"
	"github.com/riskibarqy/bq-account-service/databases"
	"github.com/riskibarqy/bq-account-service/external/clerk"
	"github.com/riskibarqy/bq-account-service/external/logger"
	"github.com/riskibarqy/bq-account-service/internal/data"
	internalhttp "github.com/riskibarqy/bq-account-service/internal/http"
	"github.com/riskibarqy/bq-account-service/internal/redis"
	"github.com/riskibarqy/bq-account-service/internal/repository/models"
	userPg "github.com/riskibarqy/bq-account-service/internal/repository/user"
	"github.com/riskibarqy/bq-account-service/internal/usecase/user"
)

var ctx = context.Background()

// InternalServices represents all the internal domain services
type InternalServices struct {
	userService user.ServiceInterface
}

func buildInternalServices(db *sqlx.DB, _ *config.Config) *InternalServices {
	userPostgresStorage := userPg.NewUserRepository(
		data.NewPostgresStorage(db, "user", models.User{}),
	)

	userService := user.NewUserService(userPostgresStorage)
	return &InternalServices{
		userService: userService,
	}
}

func initMetadataConfig() {
	rc, err := reload.New(ctx)
	if err != nil {
		log.Fatalln(err)
		return
	}

	config.MetadataConfig = &config.Metadata{}

	go func() {
		for {
			select {
			case err := <-rc.GetErrChannel():
				log.Printf("Received err: %v", err)
			case conf := <-rc.GetReloadChan():
				log.Println("Received new config [", conf.FilePath, "]:", conf.Config)
			}
		}
	}()

	err = rc.AddConfiguration("./metadata.json", &config.MetadataConfig)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
}

func main() {
	go initMetadataConfig()
	config.GetConfiguration()

	databases.Init()
	defer func() {
		if config.AppConfig.DatabaseClient != nil {
			_ = config.AppConfig.DatabaseClient.Close()
			log.Println("Database connection closed")
		}
	}()

	redis.Init()
	clerk.Init()
	logger.Init()
	defer logger.Shutdown(ctx)

	// Print the current mode
	fmt.Printf("Running in %s mode\n", config.AppConfig.AppMode)

	dataManager := data.NewManager(config.AppConfig.DatabaseClient)
	internalServices := buildInternalServices(config.AppConfig.DatabaseClient, config.AppConfig)

	s := internalhttp.NewServer(
		config.AppConfig,
		dataManager,
		internalServices.userService,
	)

	s.Serve()
}

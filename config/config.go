package config

import (
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const (
	appMode            = "APP_MODE"
	appName            = "APP_NAME"
	appVersion         = "APP_VERSION"
	appDescription     = "APP_DESCRIPTION"
	appPort            = "APP_PORT"
	dbConnectionString = "DB_CONNECTION_STRING"
	dbName             = "DB_NAME"
	jwtSecret          = "JWT_SECRET"
	clerkSecretKey     = "CLERK_SECRET_KEY"
	redisURL           = "REDIS_URL"
	uptraceDSN         = "UPTRACE_DSN"

	redisExpirationShort  = "REDIS_EXPIRATION_SHORT"
	redisExpirationMedium = "REDIS_EXPIRATION_MEDIUM"
	redisExpirationLong   = "REDIS_EXPIRATION_LONG"
)

// Config contains application configuration
type Config struct {
	AppMode            string `json:"appMode"`
	AppName            string `json:"appName"`
	AppVersion         string `json:"appVersion"`
	AppDescription     string `json:"appDescription"`
	AppPort            string `json:"appPort"`
	DBConnectionString string `json:"dbConnectionString"`
	DBName             string `json:"dbName"`
	JWTSecret          string `json:"jwtSecret"`
	ClerkSecretKey     string `json:"clerkSecret"`
	RedisURL           string `json:"redisUrl"`
	UptraceDSN         string `json:"uptraceDsn"`

	DatabaseClient *sqlx.DB
	RedisClient    *redis.UniversalClient
}

type Metadata struct {
	RedisExpirationShort  int `json:"redisExpirationShort"`
	RedisExpirationMedium int `json:"redisExpirationMedium"`
	RedisExpirationLong   int `json:"redisExpirationLong"`
}

var AppConfig = &Config{}
var MetadataConfig = &Metadata{}

// getEnvOrDefault retrieves the value of an environment variable or returns a default value.
// It supports string, int, bool, and float64 types.
func getEnvOrDefault(env string, defaultVal interface{}) interface{} {
	e, _ := os.LookupEnv(env)
	if e == "" {
		return defaultVal
	}

	switch v := defaultVal.(type) {
	case string:
		return e
	case int:
		if intVal, err := strconv.Atoi(e); err == nil {
			return intVal
		}
		return v // return default if conversion fails
	case bool:
		if boolVal, err := strconv.ParseBool(e); err == nil {
			return boolVal
		}
		return v // return default if conversion fails
	case float64:
		if floatVal, err := strconv.ParseFloat(e, 64); err == nil {
			return floatVal
		}
		return v // return default if conversion fails
	default:
		return defaultVal // return default for unsupported types
	}
}

func GetConfiguration() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	AppConfig.AppMode = getEnvOrDefault(appMode, "development").(string)
	AppConfig.AppName = getEnvOrDefault(appName, "account-service").(string)
	AppConfig.AppVersion = getEnvOrDefault(appVersion, "v1.0.0").(string)
	AppConfig.AppDescription = getEnvOrDefault(appDescription, "Account Service").(string)
	AppConfig.AppPort = getEnvOrDefault(appPort, "8080").(string)

	AppConfig.DBConnectionString = getEnvOrDefault(dbConnectionString, "postgres://user:password@localhost:5432/account_db?sslmode=disable").(string)
	AppConfig.DBName = getEnvOrDefault(dbName, "dbname").(string)
	AppConfig.JWTSecret = getEnvOrDefault(jwtSecret, "supersecret").(string)
	AppConfig.ClerkSecretKey = getEnvOrDefault(clerkSecretKey, "test").(string)
	AppConfig.RedisURL = getEnvOrDefault(redisURL, "redis://localhost:6379").(string)
	AppConfig.UptraceDSN = getEnvOrDefault(uptraceDSN, "").(string)

	MetadataConfig.RedisExpirationShort = getEnvOrDefault(redisExpirationShort, 60).(int)
	MetadataConfig.RedisExpirationMedium = getEnvOrDefault(redisExpirationMedium, 3600).(int)
	MetadataConfig.RedisExpirationLong = getEnvOrDefault(redisExpirationLong, 86400).(int)
}

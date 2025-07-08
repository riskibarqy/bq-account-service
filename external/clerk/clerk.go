package clerk

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/riskibarqy/bq-account-service/config"
)

func Init() {
	clerk.SetKey(config.AppConfig.ClerkSecretKey)
}

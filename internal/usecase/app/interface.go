package app

import (
	"context"

	"github.com/riskibarqy/bq-account-service/internal/domain/entity"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// ServiceInterface represents the app service interface
type AppServiceInterface interface {
	ListApps(ctx context.Context, params *datatransfers.FindAllParams) ([]*entity.App, int, *types.Error)
	// GetApp(ctx context.Context, appID int) (*models.App, *types.Error)
	// CreateApp(ctx context.Context, params *datatransfers.RegisterApp) (*models.App, *types.Error)
	// UpdateApp(ctx context.Context, appID int, params *models.App) (*models.App, *types.Error)
	// DeleteApp(ctx context.Context, appID int) *types.Error
	// ChangePassword(ctx context.Context, appID int, oldPassword, newPassword string) *types.Error
	// Login(ctx context.Context, email string, password string) (*datatransfers.LoginResponse, *types.Error)
}

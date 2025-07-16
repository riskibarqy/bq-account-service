package app

import (
	"context"

	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/models"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// Storage represents the app storage interface
type Storage interface {
	FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*models.App, *types.Error)
	FindByID(ctx context.Context, appID int) (*models.App, *types.Error)
	Insert(ctx context.Context, app *models.App) (*models.App, *types.Error)
	Update(ctx context.Context, app *models.App) (*models.App, *types.Error)
	Delete(ctx context.Context, appID int) *types.Error
}

package app

import (
	"context"

	"github.com/riskibarqy/bq-account-service/internal/domain/entity"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/repository/models"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// Storage represents the app storage interface
type Storage interface {
	FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*entity.App, *types.Error)
	FindByID(ctx context.Context, appID int) (*entity.App, *types.Error)
	Insert(ctx context.Context, app *models.App) (*entity.App, *types.Error)
	Update(ctx context.Context, app *models.App) (*entity.App, *types.Error)
	Delete(ctx context.Context, appID int) *types.Error
}

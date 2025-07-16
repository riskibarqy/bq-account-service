package user

import (
	"context"

	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/models"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// Storage represents the user storage interface
type Storage interface {
	FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*models.User, *types.Error)
	FindByID(ctx context.Context, userID int) (*models.User, *types.Error)
	FindByEmail(ctx context.Context, email string) (*models.User, *types.Error)
	Insert(ctx context.Context, user *models.User) (*models.User, *types.Error)
	Update(ctx context.Context, user *models.User) (*models.User, *types.Error)
	Delete(ctx context.Context, userID int) *types.Error
}

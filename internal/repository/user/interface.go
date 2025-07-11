package user

import (
	"context"

	"github.com/riskibarqy/bq-account-service/internal/domain/entity"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/repository/models"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// Storage represents the user storage interface
type Storage interface {
	FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*entity.User, *types.Error)
	FindByID(ctx context.Context, userID int) (*entity.User, *types.Error)
	FindByEmail(ctx context.Context, email string) (*entity.User, *types.Error)
	Insert(ctx context.Context, user *models.User) (*entity.User, *types.Error)
	Update(ctx context.Context, user *models.User) (*entity.User, *types.Error)
	Delete(ctx context.Context, userID int) *types.Error
}

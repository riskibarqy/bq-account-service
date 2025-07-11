package user

import (
	"context"

	"github.com/riskibarqy/bq-account-service/internal/domain/entity"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// ServiceInterface represents the user service interface
type ServiceInterface interface {
	ListUsers(ctx context.Context, params *datatransfers.FindAllParams) ([]*entity.User, int, *types.Error)
	// GetUser(ctx context.Context, userID int) (*models.User, *types.Error)
	// CreateUser(ctx context.Context, params *datatransfers.RegisterUser) (*models.User, *types.Error)
	Register(ctx context.Context, params *datatransfers.RegisterUser) (*entity.User, *types.Error)
	// UpdateUser(ctx context.Context, userID int, params *models.User) (*models.User, *types.Error)
	// DeleteUser(ctx context.Context, userID int) *types.Error
	// ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) *types.Error
	// Login(ctx context.Context, email string, password string) (*datatransfers.LoginResponse, *types.Error)
}

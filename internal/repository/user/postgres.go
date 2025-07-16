package user

import (
	"context"
	"fmt"

	"github.com/riskibarqy/bq-account-service/internal/data"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/models"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// UserRepository implements the user storage service interface
type UserRepository struct {
	Storage data.GenericStorage
}

// FindAll find all users
func (s *UserRepository) FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*models.User, *types.Error) {

	users := []*models.User{}
	where := `"deleted_at" IS NULL`

	if params.Email != "" {
		where += ` AND "email" ILIKE :email`
	}

	if params.Phone != "" {
		where += ` AND "phone" ILIKE :phone`
	}

	if params.UserID != 0 {
		where += ` AND "id" = :userId`
	}
	if params.Name != "" {
		where += ` AND "name" ILIKE :name`
	}
	if len(params.UserIDs) > 0 {
		where += ` AND "id" in (:userIds)`
	}
	if params.Page != 0 && params.Limit != 0 {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC LIMIT :limit OFFSET :offset`, where)
	} else {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC`, where)
	}

	err := s.Storage.Where(ctx, &users, where, map[string]interface{}{
		"userId":  params.UserID,
		"userIds": params.UserIDs,
		"limit":   params.Limit,
		"email":   params.Email,
		"phone":   params.Phone,
		"name":    params.Name,
		"offset":  ((params.Page - 1) * params.Limit),
	})
	if err != nil {
		return nil, types.NewError(err)
	}

	return users, nil
}

// FindByID find user by its id
func (s *UserRepository) FindByID(ctx context.Context, userID int) (*models.User, *types.Error) {
	user := &models.User{}
	err := s.Storage.FindByID(ctx, user, userID)
	if err != nil {
		return nil, types.NewError(err)
	}

	return user, nil
}

// FindByEmail find user by its email
func (s *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, *types.Error) {
	users, err := s.FindAll(ctx, &datatransfers.FindAllParams{
		Email: email,
	})
	if err != nil {
		err.Path = ".UserRepository->FindByEmail()" + err.Path
		return nil, err
	}

	if len(users) < 1 || users[0].Email != email {
		return nil, types.NewError(types.ErrNotFound)
	}

	return users[0], nil
}

// Insert insert user
func (s *UserRepository) Insert(ctx context.Context, user *models.User) (*models.User, *types.Error) {
	err := s.Storage.Insert(ctx, user)
	if err != nil {
		return nil, types.NewError(err)
	}

	return user, nil
}

// Update update user
func (s *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, *types.Error) {
	err := s.Storage.Update(ctx, user)
	if err != nil {
		return nil, types.NewError(err)
	}

	return user, nil
}

// Delete delete a user
func (s *UserRepository) Delete(ctx context.Context, userID int) *types.Error {
	err := s.Storage.Delete(ctx, userID)
	if err != nil {
		return types.NewError(err)
	}

	return nil
}

// NewUserRepository creates new user repository service
func NewUserRepository(
	storage data.GenericStorage,
) *UserRepository {
	return &UserRepository{
		Storage: storage,
	}
}

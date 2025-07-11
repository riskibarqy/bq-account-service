package user

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/riskibarqy/bq-account-service/internal/data"
	"github.com/riskibarqy/bq-account-service/internal/domain/entity"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/repository/models"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// UserRepository implements the user storage service interface
type UserRepository struct {
	Storage data.GenericStorage
}

// FindAll find all users
func (s *UserRepository) FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*entity.User, *types.Error) {

	users := []*models.User{}
	where := `"deleted_at" IS NULL`

	if params.UserID != 0 {
		where += ` AND "id" = :userId`
	}
	if params.Email != "" {
		where += ` AND "email" ILIKE :email`
	}
	if params.Name != "" {
		where += ` AND "name" ILIKE :name`
	}
	if params.Search != "" {
		where += ` AND "name" ILIKE :search`
	}
	if params.Token != "" {
		where += ` AND "token" ILIKE :token`
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
		"name":    params.Name,
		"search":  "%" + params.Search + "%",
		"offset":  ((params.Page - 1) * params.Limit),
		"token":   params.Token,
	})
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserRepository->FindAll()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	// Map models.user to entity.user
	result := make([]*entity.User, 0, len(users))
	for _, m := range users {
		e := &entity.User{}
		if err := copier.Copy(e, m); err != nil {
			return nil, &types.Error{
				Path:    ".UserAppRepository->FindAll()->copier",
				Message: err.Error(),
				Error:   err,
				Type:    "mapping-error",
			}
		}
		result = append(result, e)
	}

	return result, nil
}

// FindByID find user by its id
func (s *UserRepository) FindByID(ctx context.Context, userID int) (*entity.User, *types.Error) {
	users, err := s.FindAll(ctx, &datatransfers.FindAllParams{
		UserID: userID,
	})
	if err != nil {
		err.Path = ".UserRepository->FindByID()" + err.Path
		return nil, err
	}

	if len(users) < 1 || users[0].ID != userID {
		return nil, &types.Error{
			Path:    ".UserRepository->FindByID()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return users[0], nil
}

// FindByEmail find user by its email
func (s *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, *types.Error) {
	users, err := s.FindAll(ctx, &datatransfers.FindAllParams{
		Email: email,
	})
	if err != nil {
		err.Path = ".UserRepository->FindByEmail()" + err.Path
		return nil, err
	}

	if len(users) < 1 || users[0].Email != email {
		return nil, &types.Error{
			Path:    ".UserRepository->FindByEmail()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return users[0], nil
}

// // FindByToken find user by its token
// func (s *UserRepository) FindByToken(ctx context.Context, token string) (*models.User, *types.Error) {
// 	users, err := s.FindAll(ctx, &datatransfers.FindAllParams{
// 		Token: token,
// 	})
// 	if err != nil {
// 		err.Path = ".UserRepository->FindByToken()" + err.Path
// 		return nil, err
// 	}

// 	if len(users) < 1 || (users[0].Token != nil && *users[0].Token != token) {
// 		return nil, &types.Error{
// 			Path:    ".UserRepository->FindByToken()",
// 			Message: data.ErrNotFound.Error(),
// 			Error:   data.ErrNotFound,
// 			Type:    "pq-error",
// 		}
// 	}

// 	return users[0], nil
// }

// Insert insert user
func (s *UserRepository) Insert(ctx context.Context, user *models.User) (*entity.User, *types.Error) {
	err := s.Storage.Insert(ctx, user)
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserRepository->Insert()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return &entity.User{
		ID:         user.ID,
		ClerkID:    user.ClerkID,
		Name:       user.Name,
		Username:   user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		DeletedAt:  user.DeletedAt,
	}, nil
}

// Update update user
func (s *UserRepository) Update(ctx context.Context, user *models.User) (*entity.User, *types.Error) {
	err := s.Storage.Update(ctx, user)
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserRepository->Update()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return &entity.User{
		ID:         user.ID,
		ClerkID:    user.ClerkID,
		Name:       user.Name,
		Username:   user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		DeletedAt:  user.DeletedAt,
	}, nil
}

// Delete delete a user
func (s *UserRepository) Delete(ctx context.Context, userID int) *types.Error {
	err := s.Storage.Delete(ctx, userID)
	if err != nil {
		return &types.Error{
			Path:    ".UserRepository->Delete()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
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

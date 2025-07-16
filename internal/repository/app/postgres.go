package app

import (
	"context"

	"github.com/riskibarqy/bq-account-service/internal/data"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/models"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// AppRepository implements the app storage service interface
type AppRepository struct {
	Storage data.GenericStorage
}

// FindAll finds all apps and maps from models to entity
func (s *AppRepository) FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*models.App, *types.Error) {
	apps := []*models.App{}
	where := `"deleted_at" IS NULL`

	if params.AppID != 0 {
		where += ` AND "id" = :appId`
	}
	if params.Email != "" {
		where += ` AND "email" ILIKE :email`
	}
	if params.Name != "" {
		where += ` AND "name" ILIKE :name`
	}
	if len(params.AppIDs) > 0 {
		where += ` AND "id" in (:appIds)`
	}

	if params.Page != 0 && params.Limit != 0 {
		where += ` ORDER BY "id" DESC LIMIT :limit OFFSET :offset`
	} else {
		where += ` ORDER BY "id" DESC`
	}

	queryParams := map[string]interface{}{
		"appId":  params.AppID,
		"appIds": params.AppIDs,
		"limit":  params.Limit,
		"email":  params.Email,
		"name":   params.Name,
		"offset": (params.Page - 1) * params.Limit,
	}

	err := s.Storage.Where(ctx, &apps, where, queryParams)
	if err != nil {
		return nil, types.NewError(err)
	}

	return apps, nil
}

// FindByID find app by its id
func (s *AppRepository) FindByID(ctx context.Context, appID int) (*models.App, *types.Error) {
	app := &models.App{}
	err := s.Storage.FindByID(ctx, app, appID)
	if err != nil {
		return nil, types.NewError(err)
	}

	return app, nil
}

// Insert insert app
func (s *AppRepository) Insert(ctx context.Context, app *models.App) (*models.App, *types.Error) {
	err := s.Storage.Insert(ctx, app)
	if err != nil {
		return nil, &types.Error{
			Path:    ".AppRepository->Insert()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}
	return app, nil
}

// Update update app
func (s *AppRepository) Update(ctx context.Context, app *models.App) (*models.App, *types.Error) {
	err := s.Storage.Update(ctx, app)
	if err != nil {
		return nil, &types.Error{
			Path:    ".AppRepository->Update()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return app, nil
}

// Delete delete a app
func (s *AppRepository) Delete(ctx context.Context, appID int) *types.Error {
	err := s.Storage.Delete(ctx, appID)
	if err != nil {
		return &types.Error{
			Path:    ".AppRepository->Delete()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return nil
}

// NewAppRepository creates new app repository service
func NewAppRepository(
	storage data.GenericStorage,
) *AppRepository {
	return &AppRepository{
		Storage: storage,
	}
}

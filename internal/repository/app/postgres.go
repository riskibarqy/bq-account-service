package app

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/riskibarqy/bq-account-service/internal/data"
	"github.com/riskibarqy/bq-account-service/internal/domain/entity"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/repository/models"
	"github.com/riskibarqy/bq-account-service/internal/types"
)

// AppRepository implements the app storage service interface
type AppRepository struct {
	Storage data.GenericStorage
}

// FindAll finds all apps and maps from models to entity
func (s *AppRepository) FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*entity.App, *types.Error) {
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
	if params.Search != "" {
		where += ` AND "name" ILIKE :search`
	}
	if params.Token != "" {
		where += ` AND "token" ILIKE :token`
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
		"search": "%" + params.Search + "%",
		"offset": (params.Page - 1) * params.Limit,
		"token":  params.Token,
	}

	err := s.Storage.Where(ctx, &apps, where, queryParams)
	if err != nil {
		return nil, &types.Error{
			Path:    ".AppRepository->FindAll()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	// Map models.App to entity.App
	result := make([]*entity.App, 0, len(apps))
	for _, m := range apps {
		e := &entity.App{}
		if err := copier.Copy(e, m); err != nil {
			return nil, &types.Error{
				Path:    ".AppRepository->FindAll()->copier",
				Message: err.Error(),
				Error:   err,
				Type:    "mapping-error",
			}
		}
		result = append(result, e)
	}

	return result, nil
}

// FindByID find app by its id
func (s *AppRepository) FindByID(ctx context.Context, appID int) (*entity.App, *types.Error) {
	apps, err := s.FindAll(ctx, &datatransfers.FindAllParams{
		AppID: appID,
	})
	if err != nil {
		err.Path = ".AppAppRepository->FindByID()" + err.Path
		return nil, err
	}

	if len(apps) < 1 || apps[0].ID != appID {
		return nil, &types.Error{
			Path:    ".AppAppRepository->FindByID()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return apps[0], nil
}

// Insert insert app
func (s *AppRepository) Insert(ctx context.Context, app *models.App) (*entity.App, *types.Error) {
	err := s.Storage.Insert(ctx, app)
	if err != nil {
		return nil, &types.Error{
			Path:    ".AppAppRepository->Insert()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return &entity.App{
		ID:           app.ID,
		Name:         app.Name,
		Slug:         app.Slug,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		CreatedAt:    app.CreatedAt,
		UpdatedAt:    &app.CreatedAt,
		DeletedAt:    app.DeletedAt,
	}, nil
}

// Update update app
func (s *AppRepository) Update(ctx context.Context, app *models.App) (*entity.App, *types.Error) {
	err := s.Storage.Update(ctx, app)
	if err != nil {
		return nil, &types.Error{
			Path:    ".AppAppRepository->Update()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return &entity.App{
		ID:           app.ID,
		Name:         app.Name,
		Slug:         app.Slug,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		CreatedAt:    app.CreatedAt,
		UpdatedAt:    &app.CreatedAt,
		DeletedAt:    app.DeletedAt,
	}, nil
}

// Delete delete a app
func (s *AppRepository) Delete(ctx context.Context, appID int) *types.Error {
	err := s.Storage.Delete(ctx, appID)
	if err != nil {
		return &types.Error{
			Path:    ".AppAppRepository->Delete()",
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

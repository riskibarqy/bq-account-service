package app

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/riskibarqy/bq-account-service/config"
	"github.com/riskibarqy/bq-account-service/internal/domain/entity"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/redis"
	"github.com/riskibarqy/bq-account-service/internal/repository/app"
	"github.com/riskibarqy/bq-account-service/internal/types"
	"github.com/riskibarqy/bq-account-service/utils"
)

// Service is the domain logic implementation of app Service interface
type AppService struct {
	appStorage app.Storage
}

func (s *AppService) ListApps(ctx context.Context, params *datatransfers.FindAllParams) ([]*entity.App, int, *types.Error) {
	// Generate cache key
	byteParams, _ := jsoniter.Marshal(params)
	cacheKey := fmt.Sprintf("ListApps-%s", utils.EncodeHexMD5(string(byteParams)))

	// Try to get apps from Redis cache
	cached, count, errCache := redis.GetListCache(ctx, cacheKey)
	if errCache == nil && cached != "" {
		// If cache hit, unmarshal the cached data into a slice of App entity
		var apps []*entity.App
		if err := jsoniter.Unmarshal([]byte(cached), &apps); err == nil {
			return apps, count, nil
		}
	}

	// Fetch apps from database
	apps, err := s.appStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".AppService->ListApps()" + err.Path
		return nil, 0, err
	}

	go func() {
		ctxChild := context.Background()

		// Cache apps and their count
		byteResults, _ := jsoniter.Marshal(apps)
		expiration := time.Duration(config.MetadataConfig.RedisExpirationShort) * time.Second

		if err := redis.SetCache(ctxChild, cacheKey, byteResults, expiration); err != nil {
			log.Printf("Failed to set app cache: %v", err)
		}

		if err := redis.SetCache(ctxChild, fmt.Sprintf("cnt-%s", cacheKey), strconv.Itoa(len(apps)), expiration); err != nil {
			log.Printf("Failed to set app count cache: %v", err)
		}
	}()

	return apps, len(apps), nil
}

// // GetApp is get app
// func (s *Service) GetApp(ctx context.Context, appID int) (*App, *types.Error) {
// 	cacheKey := fmt.Sprintf("GetApp-%d", appID)

// 	// Try to get apps from Redis cache
// 	cached, errCache := redis.GetCache(ctx, cacheKey)
// 	if errCache == nil && cached != "" {
// 		// If cache hit, unmarshal the cached data into a slice of App entity
// 		var app *App
// 		if err := jsoniter.Unmarshal([]byte(cached), &app); err == nil {
// 			return app, nil
// 		}
// 	}

// 	app, err := s.appStorage.FindByID(ctx, appID)
// 	if err != nil {
// 		err.Path = ".AppService->GetApp()" + err.Path
// 		return nil, err
// 	}

// 	go func() {
// 		ctxChild := context.Background()

// 		// Cache app
// 		byteResults, _ := jsoniter.Marshal(app)
// 		expiration := time.Duration(config.MetadataConfig.RedisExpirationShort) * time.Second

// 		if err := redis.SetCache(ctxChild, cacheKey, byteResults, expiration); err != nil {
// 			log.Printf("Failed to set app cache: %v", err)
// 		}
// 	}()

// 	return app, nil
// }

// // CreateApp create app
// func (s *Service) CreateApp(ctx context.Context, params *App) (*App, *types.Error) {
// 	apps, _, errType := s.ListApps(ctx, &datatransfers.FindAllParams{
// 		Email: params.Email,
// 	})
// 	if errType != nil {
// 		errType.Path = ".AppService->CreateApp()" + errType.Path
// 		return nil, errType
// 	}
// 	if len(apps) > 0 {
// 		return nil, &types.Error{
// 			Path:    ".AppService->CreateApp()",
// 			Message: ErrEmailAlreadyExists.Error(),
// 			Error:   ErrEmailAlreadyExists,
// 			Type:    "validation-error",
// 		}
// 	}

// 	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, &types.Error{
// 			Path:    ".AppService->CreateApp()",
// 			Message: err.Error(),
// 			Error:   err,
// 			Type:    "golang-error",
// 		}
// 	}

// 	now := utils.Now()

// 	app := &App{
// 		Name:     params.Name,
// 		Email:    params.Email,
// 		Password: string(bcryptHash),
// 		// Token:          nil,
// 		// TokenExpiredAt: nil,
// 		CreatedAt: now,
// 		UpdatedAt: &now,
// 	}

// 	app, errType = s.appStorage.Insert(ctx, app)
// 	if errType != nil {
// 		errType.Path = ".AppService->CreateApp()" + errType.Path
// 		return nil, errType
// 	}

// 	return app, nil
// }

// // UpdateApp update a app
// func (s *Service) UpdateApp(ctx context.Context, appID int, params *App) (*App, *types.Error) {
// 	app, err := s.GetApp(ctx, appID)
// 	if err != nil {
// 		err.Path = ".AppService->UpdateApp()" + err.Path
// 		return nil, err
// 	}

// 	apps, _, err := s.ListApps(ctx, &datatransfers.FindAllParams{
// 		Email: params.Email,
// 	})
// 	if err != nil {
// 		err.Path = ".AppService->UpdateApp()" + err.Path
// 		return nil, err
// 	}
// 	if len(apps) > 0 {
// 		return nil, &types.Error{
// 			Path:    ".AppService->CreateApp()",
// 			Message: data.ErrAlreadyExist.Error(),
// 			Error:   data.ErrAlreadyExist,
// 			Type:    "validation-error",
// 		}
// 	}

// 	app.Name = params.Name
// 	app.Email = params.Email

// 	app, err = s.appStorage.Update(ctx, app)
// 	if err != nil {
// 		err.Path = ".AppService->UpdateApp()" + err.Path
// 		return nil, err
// 	}

// 	go func() {
// 		ctxChild := context.Background()

// 		cacheKey := fmt.Sprintf("GetApp-%d", appID)

// 		// delete app cache
// 		if err := redis.DeleteCache(ctxChild, cacheKey); err != nil {
// 			log.Printf("Failed to set app cache: %v", err)
// 		}
// 	}()

// 	return app, nil
// }

// // DeleteApp delete a app
// func (s *Service) DeleteApp(ctx context.Context, appID int) *types.Error {
// 	err := s.appStorage.Delete(ctx, appID)
// 	if err != nil {
// 		err.Path = ".AppService->DeleteApp()" + err.Path
// 		return err
// 	}

// 	go func() {
// 		ctxChild := context.Background()

// 		cacheKey := fmt.Sprintf("GetApp-%d", appID)

// 		// delete app cache
// 		if err := redis.DeleteCache(ctxChild, cacheKey); err != nil {
// 			log.Printf("Failed to set app cache: %v", err)
// 		}
// 	}()

// 	return nil
// }

// // ChangePassword change password
// func (s *Service) ChangePassword(ctx context.Context, appID int, oldPassword, newPassword string) *types.Error {
// 	app, err := s.GetApp(ctx, appID)
// 	if err != nil {
// 		err.Path = ".AppService->ChangePassword()" + err.Path
// 		return err
// 	}

// 	errBcrypt := bcrypt.CompareHashAndPassword([]byte(app.Password), []byte(oldPassword))
// 	if errBcrypt != nil {
// 		return &types.Error{
// 			Path:    ".AppService->ChangePassword()",
// 			Message: ErrWrongPassword.Error(),
// 			Error:   ErrWrongPassword,
// 			Type:    "golang-error",
// 		}
// 	}

// 	bcryptHash, errBcrypt := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 	if errBcrypt != nil {
// 		return &types.Error{
// 			Path:    ".AppService->ChangePassword()",
// 			Message: errBcrypt.Error(),
// 			Error:   errBcrypt,
// 			Type:    "golang-error",
// 		}
// 	}

// 	app.Password = string(bcryptHash)
// 	_, err = s.appStorage.Update(ctx, app)
// 	if err != nil {
// 		err.Path = ".AppService->ChangePassword()" + err.Path
// 		return err
// 	}

// 	go func() {
// 		ctxChild := context.Background()

// 		cacheKey := fmt.Sprintf("GetApp-%d", appID)

// 		// delete app cache
// 		if err := redis.DeleteCache(ctxChild, cacheKey); err != nil {
// 			log.Printf("Failed to set app cache: %v", err)
// 		}
// 	}()

// 	return nil
// }

// // Login login
// func (s *Service) Login(ctx context.Context, email string, password string) (*datatransfers.LoginResponse, *types.Error) {
// 	apps, _, err := s.ListApps(ctx, &datatransfers.FindAllParams{
// 		Email: email,
// 	})
// 	if err != nil {
// 		err.Path = ".AppService->Login()" + err.Path
// 		return nil, err
// 	}
// 	if len(apps) < 1 {
// 		return nil, &types.Error{
// 			Path:    ".AppService->Login()",
// 			Message: ErrWrongEmail.Error(),
// 			Error:   ErrWrongEmail,
// 			Type:    "validation-error",
// 		}
// 	}

// 	app := apps[0]
// 	errBcrypt := bcrypt.CompareHashAndPassword([]byte(app.Password), []byte(password))
// 	if errBcrypt != nil {
// 		return nil, &types.Error{
// 			Path:    ".AppService->ChangePassword()",
// 			Message: ErrWrongPassword.Error(),
// 			Error:   ErrWrongPassword,
// 			Type:    "golang-error",
// 		}
// 	}

// 	token, errToken := config.GenerateJWTToken(app)
// 	if errToken != nil {
// 		return nil, &types.Error{
// 			Path:    ".AppService->CreateApp()",
// 			Message: errToken.Error(),
// 			Error:   errToken,
// 			Type:    "golang-error",
// 		}
// 	}

// 	now := utils.Now()
// 	// tokenExpiredAt := now + 72*3600

// 	// app.Token = &token
// 	// app.TokenExpiredAt = &tokenExpiredAt
// 	app.UpdatedAt = &now

// 	app, err = s.appStorage.Update(ctx, app)
// 	if err != nil {
// 		err.Path = ".AppService->CreateApp()" + err.Path
// 		return nil, err
// 	}

// 	return &datatransfers.LoginResponse{
// 		SessionID: token,
// 		App:      app,
// 	}, nil
// }

// NewService creates a new app AppService
func NewService(
	appStorage app.Storage,
) *AppService {
	return &AppService{
		appStorage: appStorage,
	}
}

package user

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	jsoniter "github.com/json-iterator/go"
	"github.com/riskibarqy/bq-account-service/config"
	"github.com/riskibarqy/bq-account-service/external/logger"
	"github.com/riskibarqy/bq-account-service/external/redis"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/models"
	"github.com/riskibarqy/bq-account-service/internal/repository/user"
	"github.com/riskibarqy/bq-account-service/internal/types"
	"github.com/riskibarqy/bq-account-service/utils"
)

// Service is the domain logic implementation of user Service interface
type Service struct {
	userStorage user.Storage
}

func (s *Service) ListUsers(ctx context.Context, params *datatransfers.FindAllParams) ([]*models.User, int, *types.Error) {
	// Generate cache key
	byteParams, _ := jsoniter.Marshal(params)
	cacheKey := fmt.Sprintf("ListUsers-%s", utils.EncodeHexMD5(string(byteParams)))

	// Try to get users from Redis cache
	cached, count, errCache := redis.GetListCache(ctx, cacheKey)
	if errCache == nil && cached != "" {
		// If cache hit, unmarshal the cached data into a slice of User models
		var users []*models.User
		if err := jsoniter.Unmarshal([]byte(cached), &users); err == nil {
			return users, count, nil
		}
	}

	// Fetch users from database
	users, err := s.userStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserService->ListUsers()" + err.Path
		return nil, 0, err
	}

	go func() {
		ctxChild := context.Background()

		// Cache users and their count
		byteResults, _ := jsoniter.Marshal(users)
		expiration := time.Duration(config.MetadataConfig.RedisExpirationShort) * time.Second

		if err := redis.SetCache(ctxChild, cacheKey, byteResults, expiration); err != nil {
			log.Printf("Failed to set user cache: %v", err)
		}

		if err := redis.SetCache(ctxChild, fmt.Sprintf("cnt-%s", cacheKey), strconv.Itoa(len(users)), expiration); err != nil {
			log.Printf("Failed to set user count cache: %v", err)
		}
	}()

	return users, len(users), nil
}

// // GetUser is get user
// func (s *Service) GetUser(ctx context.Context, userID int) (*models.User, *types.Error) {
// 	cacheKey := fmt.Sprintf("GetUser-%d", userID)

// 	// Try to get users from Redis cache
// 	cached, errCache := redis.GetCache(ctx, cacheKey)
// 	if errCache == nil && cached != "" {
// 		// If cache hit, unmarshal the cached data into a slice of User models
// 		var user *models.User
// 		if err := jsoniter.Unmarshal([]byte(cached), &user); err == nil {
// 			return user, nil
// 		}
// 	}

// 	user, err := s.userStorage.FindByID(ctx, userID)
// 	if err != nil {
// 		err.Path = ".UserService->GetUser()" + err.Path
// 		return nil, err
// 	}

// 	go func() {
// 		ctxChild := context.Background()

// 		// Cache user
// 		byteResults, _ := jsoniter.Marshal(user)
// 		expiration := time.Duration(config.MetadataConfig.RedisExpirationShort) * time.Second

// 		if err := redis.SetCache(ctxChild, cacheKey, byteResults, expiration); err != nil {
// 			log.Printf("Failed to set user cache: %v", err)
// 		}
// 	}()

// 	return user, nil
// }

// Register create user
func (s *Service) Register(ctx context.Context, params *datatransfers.RegisterUser) (*models.User, *types.Error) {
	users, _, errType := s.ListUsers(ctx, &datatransfers.FindAllParams{
		Email: params.Email,
		Phone: params.Phone,
	})
	if errType != nil {
		return nil, errType
	}

	if len(users) > 0 {
		return nil, types.NewError(types.ErrUserAlreadyExists)
	}

	f, l := utils.SplitName(params.Name)
	if params.Username == "" {
		params.Username = utils.CreateUsernameFromEmail(params.Email)
	}

	clerkCreateResponse, errClerk := clerkUser.Create(ctx, &clerkUser.CreateParams{
		EmailAddresses: &[]string{params.Email},
		Username:       &params.Username,
		Password:       &params.Password,
		FirstName:      &f,
		LastName:       &l,
	})
	if errClerk != nil {
		return nil, &types.Error{
			Path:    ".UserService->Register()",
			Message: errClerk.Error(),
			Error:   errClerk,
			Type:    "clerk-create-user",
		}
	}

	now := utils.Now()
	userModel := &models.User{
		ClerkID:   clerkCreateResponse.ID,
		Name:      *clerkCreateResponse.FirstName + " " + *clerkCreateResponse.LastName,
		Email:     params.Email,
		Username:  *clerkCreateResponse.Username,
		Phone:     params.Phone,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: &now,
	}

	user, errType := s.userStorage.Insert(ctx, userModel)
	if errType != nil {
		ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if _, errClerkDeleteUser := clerkUser.Delete(ctxTimeout, clerkCreateResponse.ID); errClerkDeleteUser != nil {
			(&types.Error{
				Path:    ".UserService->Register()",
				Message: errClerkDeleteUser.Error(),
				Error:   errClerkDeleteUser,
				Type:    "clerk-delete-user",
			}).Log(ctx, logger.Tracer)
		}

		errType.Path = ".UserService->Register()" + errType.Path
		return nil, errType
	}

	return user, nil
}

// // UpdateUser update a user
// func (s *Service) UpdateUser(ctx context.Context, userID int, params *models.User) (*models.User, *types.Error) {
// 	user, err := s.GetUser(ctx, userID)
// 	if err != nil {
// 		err.Path = ".UserService->UpdateUser()" + err.Path
// 		return nil, err
// 	}

// 	users, _, err := s.ListUsers(ctx, &datatransfers.FindAllParams{
// 		Email: params.Email,
// 	})
// 	if err != nil {
// 		err.Path = ".UserService->UpdateUser()" + err.Path
// 		return nil, err
// 	}
// 	if len(users) > 0 {
// 		return nil, &types.Error{
// 			Path:    ".UserService->CreateUser()",
// 			Message: data.ErrAlreadyExist.Error(),
// 			Error:   data.ErrAlreadyExist,
// 			Type:    "validation-error",
// 		}
// 	}

// 	user.Name = params.Name
// 	user.Email = params.Email

// 	user, err = s.userStorage.Update(ctx, user)
// 	if err != nil {
// 		err.Path = ".UserService->UpdateUser()" + err.Path
// 		return nil, err
// 	}

// 	go func() {
// 		ctxChild := context.Background()

// 		cacheKey := fmt.Sprintf("GetUser-%d", userID)

// 		// delete user cache
// 		if err := redis.DeleteCache(ctxChild, cacheKey); err != nil {
// 			log.Printf("Failed to set user cache: %v", err)
// 		}
// 	}()

// 	return user, nil
// }

// // DeleteUser delete a user
// func (s *Service) DeleteUser(ctx context.Context, userID int) *types.Error {
// 	err := s.userStorage.Delete(ctx, userID)
// 	if err != nil {
// 		err.Path = ".UserService->DeleteUser()" + err.Path
// 		return err
// 	}

// 	go func() {
// 		ctxChild := context.Background()

// 		cacheKey := fmt.Sprintf("GetUser-%d", userID)

// 		// delete user cache
// 		if err := redis.DeleteCache(ctxChild, cacheKey); err != nil {
// 			log.Printf("Failed to set user cache: %v", err)
// 		}
// 	}()

// 	return nil
// }

// // ChangePassword change password
// func (s *Service) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) *types.Error {
// 	user, err := s.GetUser(ctx, userID)
// 	if err != nil {
// 		err.Path = ".UserService->ChangePassword()" + err.Path
// 		return err
// 	}

// 	errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
// 	if errBcrypt != nil {
// 		return &types.Error{
// 			Path:    ".UserService->ChangePassword()",
// 			Message: ErrWrongPassword.Error(),
// 			Error:   ErrWrongPassword,
// 			Type:    "golang-error",
// 		}
// 	}

// 	bcryptHash, errBcrypt := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 	if errBcrypt != nil {
// 		return &types.Error{
// 			Path:    ".UserService->ChangePassword()",
// 			Message: errBcrypt.Error(),
// 			Error:   errBcrypt,
// 			Type:    "golang-error",
// 		}
// 	}

// 	user.Password = string(bcryptHash)
// 	_, err = s.userStorage.Update(ctx, user)
// 	if err != nil {
// 		err.Path = ".UserService->ChangePassword()" + err.Path
// 		return err
// 	}

// 	go func() {
// 		ctxChild := context.Background()

// 		cacheKey := fmt.Sprintf("GetUser-%d", userID)

// 		// delete user cache
// 		if err := redis.DeleteCache(ctxChild, cacheKey); err != nil {
// 			log.Printf("Failed to set user cache: %v", err)
// 		}
// 	}()

// 	return nil
// }

// // Login login
// func (s *Service) Login(ctx context.Context, email string, password string) (*datatransfers.LoginResponse, *types.Error) {
// 	users, _, err := s.ListUsers(ctx, &datatransfers.FindAllParams{
// 		Email: email,
// 	})
// 	if err != nil {
// 		err.Path = ".UserService->Login()" + err.Path
// 		return nil, err
// 	}
// 	if len(users) < 1 {
// 		return nil, &types.Error{
// 			Path:    ".UserService->Login()",
// 			Message: ErrWrongEmail.Error(),
// 			Error:   ErrWrongEmail,
// 			Type:    "validation-error",
// 		}
// 	}

// 	user := users[0]
// 	errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
// 	if errBcrypt != nil {
// 		return nil, &types.Error{
// 			Path:    ".UserService->ChangePassword()",
// 			Message: ErrWrongPassword.Error(),
// 			Error:   ErrWrongPassword,
// 			Type:    "golang-error",
// 		}
// 	}

// 	token, errToken := config.GenerateJWTToken(user)
// 	if errToken != nil {
// 		return nil, &types.Error{
// 			Path:    ".UserService->CreateUser()",
// 			Message: errToken.Error(),
// 			Error:   errToken,
// 			Type:    "golang-error",
// 		}
// 	}

// 	now := utils.Now()
// 	// tokenExpiredAt := now + 72*3600

// 	// user.Token = &token
// 	// user.TokenExpiredAt = &tokenExpiredAt
// 	user.UpdatedAt = &now

// 	user, err = s.userStorage.Update(ctx, user)
// 	if err != nil {
// 		err.Path = ".UserService->CreateUser()" + err.Path
// 		return nil, err
// 	}

// 	return &datatransfers.LoginResponse{
// 		SessionID: token,
// 		User:      user,
// 	}, nil
// }

// NewService creates a new user AppService
func NewUserService(
	userStorage user.Storage,
) *Service {
	return &Service{
		userStorage: userStorage,
	}
}

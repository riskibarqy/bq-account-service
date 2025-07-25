package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/riskibarqy/bq-account-service/internal/data"
	"github.com/riskibarqy/bq-account-service/internal/dto/datatransfers"
	"github.com/riskibarqy/bq-account-service/internal/http/response"
	"github.com/riskibarqy/bq-account-service/internal/models"
	"github.com/riskibarqy/bq-account-service/internal/types"
	"github.com/riskibarqy/bq-account-service/internal/usecase/user"
	"gopkg.in/go-playground/validator.v9"
)

// UserController represents the user controller
type UserController struct {
	userService user.ServiceInterface
	dataManager *data.Manager
}

// UserList user list and count
type UserList struct {
	Data  []*models.User `json:"data"`
	Count int            `json:"count"`
}

// func (a *UserController) Login(w http.ResponseWriter, r *http.Request) {
// 	var err *types.Error

// 	decoder := json.NewDecoder(r.Body)

// 	var params datatransfers.LoginParams
// 	errDecode := decoder.Decode(&params)
// 	if errDecode != nil {
// 		err = &types.Error{
// 			Path:    ".UserController->Login()",
// 			Message: errDecode.Error(),
// 			Error:   errDecode,
// 			Type:    types.ErrTypesHandlerError,
// 		}
// 		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
// 		return
// 	}

// 	var sess *datatransfers.LoginResponse
// 	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
// 		sess, err = a.userService.Login(r.Context(), params.Email, params.Password)
// 		if err != nil {
// 			return err.Error
// 		}
// 		return nil
// 	})
// 	if errTransaction != nil {
// 		err.Path = ".UserController->Login()" + err.Path
// 		if err.Error == user.ErrWrongPassword || err.Error == data.ErrNotFound {
// 			response.Error(w, "Email / password is wrong", http.StatusBadRequest, *err)
// 		} else {
// 			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
// 		}
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:  "sessionId",
// 		Value: sess.SessionID,
// 	})

// 	response.JSON(w, http.StatusOK, sess)
// }

// func (a *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
// 	var err *types.Error

// 	decoder := json.NewDecoder(r.Body)
// 	var params datatransfers.ChangePasswordParams
// 	errDecode := decoder.Decode(&params)
// 	if errDecode != nil {
// 		err = &types.Error{
// 			Path:    ".UserController->ChangePassword()",
// 			Message: errDecode.Error(),
// 			Error:   errDecode,
// 			Type:    types.ErrTypesHandlerError,
// 		}
// 		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
// 		return
// 	}

// 	userID := appcontext.UserID(r.Context())

// 	err = a.userService.ChangePassword(r.Context(), userID, params.OldPassword, params.NewPassword)
// 	if err != nil {
// 		err.Path = ".UserController->ChangePassword()" + err.Path
// 		if err.Error == user.ErrWrongPassword {
// 			response.Error(w, "Wrong old password", http.StatusBadRequest, *err)
// 		} else {
// 			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
// 		}
// 		return
// 	}

// 	response.JSON(w, http.StatusNoContent, "")
// }

// func (a *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	var err *types.Error

// 	decoder := json.NewDecoder(r.Body)

// 	var params *models.User
// 	errDecode := decoder.Decode(&params)
// 	if errDecode != nil {
// 		err = &types.Error{
// 			Path:    ".UserController->UpdateUser()",
// 			Message: errDecode.Error(),
// 			Error:   errDecode,
// 			Type:    types.ErrTypesHandlerError,
// 		}
// 		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
// 		return
// 	}
// 	var sUserID = chi.URLParam(r, "userId")
// 	userID, errConversion := strconv.Atoi(sUserID)
// 	if errConversion != nil {
// 		err = &types.Error{
// 			Path:    ".UserController->UpdateUser()",
// 			Message: errConversion.Error(),
// 			Error:   errConversion,
// 			Type:    types.ErrTypesHandlerError,
// 		}
// 		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
// 		return
// 	}

// 	var singleUser *models.User
// 	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
// 		singleUser, err = a.userService.UpdateUser(ctx, userID, params)
// 		if err != nil {
// 			return err.Error
// 		}
// 		return nil
// 	})
// 	if errTransaction != nil {
// 		err.Path = ".UserController->UpdateUser()" + err.Path
// 		if errTransaction == user.ErrEmailAlreadyExists {
// 			response.Error(w, "email has been registered", http.StatusUnprocessableEntity, *err)
// 		} else {
// 			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
// 		}
// 		return
// 	}
// 	response.JSON(w, http.StatusOK, singleUser)

// }

func (a *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var err *types.Error
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)

	var params *datatransfers.RegisterUser
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->Register()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    types.ErrTypesHandlerError,
		}
		response.Error(ctx, w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	validate := validator.New()
	errValidation := validate.Struct(params)
	if errValidation != nil {
		err = &types.Error{
			Path:    ".UserController->Register()",
			Message: errValidation.Error(),
			Error:   errValidation,
			Type:    types.ErrTypesHandlerError,
		}
		response.Error(ctx, w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	result := &models.User{}
	errTransaction := a.dataManager.RunInTransaction(ctx, func(ctx context.Context) error {
		result, err = a.userService.Register(ctx, params)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".UserController->Register()" + err.Path
		if errTransaction == types.ErrUserAlreadyExists {
			response.Error(ctx, w, types.ErrUserAlreadyExists.Error(), http.StatusUnprocessableEntity, *err)
		} else {
			response.Error(ctx, w, "Internal Server Error", http.StatusInternalServerError, *err)
		}
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// func (a *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	var err *types.Error
// 	var sUserID = chi.URLParam(r, "userId")
// 	userID, errConversion := strconv.Atoi(sUserID)
// 	if errConversion != nil {
// 		err = &types.Error{
// 			Path:    ".UserController->DeleteUser()",
// 			Message: errConversion.Error(),
// 			Error:   errConversion,
// 			Type:    types.ErrTypesHandlerError,
// 		}
// 		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
// 		return
// 	}

// 	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
// 		err = a.userService.DeleteUser(ctx, userID)
// 		if err != nil {
// 			return err.Error
// 		}
// 		return nil
// 	})
// 	if errTransaction != nil {
// 		err.Path = ".USerController->DeleteUser()" + err.Path
// 		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
// 		return
// 	}
// 	response.JSON(w, http.StatusNoContent, "")

// }

func (a *UserController) ListUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error
	ctx := r.Context()

	queryValues := r.URL.Query()
	var limit = 10
	var errConversion error
	if queryValues.Get("limit") != "" {
		limit, errConversion = strconv.Atoi(queryValues.Get("limit"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    types.ErrTypesHandlerError,
			}
			response.Error(ctx, w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var page = 1
	if queryValues.Get("page") != "" {
		page, errConversion = strconv.Atoi(queryValues.Get("page"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    types.ErrTypesHandlerError,
			}
			response.Error(ctx, w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	if limit < 0 {
		limit = 10
	}
	if page < 0 {
		page = 1
	}
	userList, count, err := a.userService.ListUsers(r.Context(), &datatransfers.FindAllParams{
		Limit: limit,
		Page:  page,
	})
	if err != nil {
		err.Path = ".UserController->ListUser()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(ctx, w, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}
	if userList == nil {
		userList = []*models.User{}
	}

	response.JSON(w, http.StatusOK, UserList{
		Data:  userList,
		Count: count,
	})
}

// func (a *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
// 	var err *types.Error

// 	var sUserID = chi.URLParam(r, "userId")
// 	userID, errConversion := strconv.Atoi(sUserID)
// 	if errConversion != nil {
// 		err = &types.Error{
// 			Path:    ".UserController->UpdateUser()",
// 			Message: errConversion.Error(),
// 			Error:   errConversion,
// 			Type:    types.ErrTypesHandlerError,
// 		}
// 		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
// 		return
// 	}

// 	user, err := a.userService.GetUser(r.Context(), userID)
// 	if err != nil {
// 		err.Path = ".UserController->GetUserByID()" + err.Path
// 		response.Error(w, "User Not Found", http.StatusNotFound, *err)
// 		if err.Error != data.ErrNotFound {
// 			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
// 			return
// 		}
// 		return
// 	}

// 	response.JSON(w, http.StatusOK, user)

// }

// NewUserController creates a new user controller
func NewUserController(
	userService user.ServiceInterface,
	dataManager *data.Manager,
) *UserController {
	return &UserController{
		userService: userService,
		dataManager: dataManager,
	}
}

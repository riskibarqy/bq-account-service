package datatransfers

import "github.com/riskibarqy/bq-account-service/models"

// LoginParams represent the http request data for login user
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response of login function
type LoginResponse struct {
	SessionID string       `json:"sessionId"`
	User      *models.User `json:"user"`
}

// ChangePasswordParams represent the http request data for change password
type ChangePasswordParams struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type RegisterUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

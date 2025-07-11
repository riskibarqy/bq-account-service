package models

// App models
type App struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name" validate:"required"`
	Slug         string `json:"slug" db:"slug" validate:"required"`
	ClientID     string `json:"clientId" db:"client_id" validate:"required"`
	ClientSecret string `json:"clientSecret" db:"client_secret" validate:"required"`
	CreatedAt    int    `json:"createdAt" db:"created_at"`
	UpdatedAt    *int   `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt    *int   `json:"deletedAt,omitempty" db:"deleted_at"`
}

func (u *App) ForPublic() {
	u.UpdatedAt = nil
	u.DeletedAt = nil
}

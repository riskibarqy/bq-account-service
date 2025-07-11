package entity

type App struct {
	ID           int
	Name         string
	Slug         string
	ClientID     string
	ClientSecret string
	CreatedAt    int
	UpdatedAt    *int
	DeletedAt    *int
}

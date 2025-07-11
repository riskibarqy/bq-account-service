package entity

type User struct {
	ID         int
	ClerkID    string
	Name       string
	Email      string
	Username   string
	Phone      string
	IsActive   bool
	IsVerified bool
	CreatedAt  int
	UpdatedAt  *int
	DeletedAt  *int
}

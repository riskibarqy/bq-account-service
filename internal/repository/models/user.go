package models

// User models
type User struct {
	ID         int    `json:"id" db:"id"`
	ClerkID    string `json:"clerkId" db:"clerk_id"`
	Name       string `json:"name" db:"name" validate:"required"`
	Email      string `json:"email" db:"email" validate:"required,email"`
	Username   string `json:"username" db:"username"`
	Phone      string `json:"phone" db:"phone"`
	IsActive   bool   `json:"isActive" db:"is_active"`
	IsVerified bool   `json:"isVerified" db:"is_verified"`
	CreatedAt  int    `json:"createdAt" db:"created_at"`
	UpdatedAt  *int   `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt  *int   `json:"deletedAt,omitempty" db:"deleted_at"`
}

func (u *User) ForPublic() {
	u.UpdatedAt = nil
	u.DeletedAt = nil
}

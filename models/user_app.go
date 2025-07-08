package models

// App models
type UserApp struct {
	ID        int         `json:"id" db:"id"`
	UserID    int         `json:"userId" db:"user_id"`
	AppID     int         `json:"appId" db:"app_id"`
	Role      int         `json:"role" db:"role"`
	Metadata  interface{} `json:"metadata" db:"metadata"`
	JoinedAt  int         `json:"joinedAt" db:"joined_at"`
	CreatedAt int         `json:"createdAt" db:"created_at"`
	UpdatedAt *int        `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt *int        `json:"deletedAt,omitempty" db:"deleted_at"`
}

func (u *UserApp) ForPublic() {
	u.UpdatedAt = nil
	u.DeletedAt = nil
}

package users

import "time"

type User struct {
	ID        int        `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Username  string     `json:"username" db:"username"`
	Password  string     `json:"-" db:"password_hash"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	LastSeen  *time.Time `json:"last_seen,omitempty" db:"last_seen"`
	IsActive  bool       `json:"is_active" db:"is_active"`
}

type UserRegister struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
	ExpiresAt   time.Time    `json:"expires_at"`
}

type UserResponse struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"created_at"`
	LastSeen  *time.Time `json:"last_seen,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	Bio       *string    `json:"bio,omitempty"`
}

type UserUpdate struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Bio       *string `json:"bio,omitempty" binding:"omitempty,max=500"`
}

type ChangePassword struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
		LastSeen:  u.LastSeen,
	}
}

package auth

import (
	"time"
)

// LoginRequest is the payload for POST /auth/login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// TokenResponse is returned after successful login or token refresh
type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	User         UserInfo `json:"user"`
}

// UserInfo is the user data included in token response and GET /me
type UserInfo struct {
	ID       uint64 `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	RoleID   uint64 `json:"role_id" db:"role_id"`
	RoleName string `json:"role_name" db:"role_name"`
	Status   string `json:"status" db:"status"`
}

// UserWithPassword is used internally for credential verification
type UserWithPassword struct {
	ID           uint64 `db:"id"`
	Name         string `db:"name"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
	RoleID       uint64 `db:"role_id"`
	RoleName     string `db:"role_name"`
	Status       string `db:"status"`
}

// RefreshRequest is the payload for POST /auth/refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordRequest is the payload for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// Claims represents the JWT claims payload
type Claims struct {
	UserID   uint64 `json:"user_id"`
	Email    string `json:"email"`
	RoleID   uint64 `json:"role_id"`
	RoleName string `json:"role_name"`
}

// Permission represents a single permission
type Permission struct {
	ID     uint64 `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Module string `db:"module" json:"module"`
}

// LoginAttempt tracks the last login time
type LoginAttempt struct {
	UserID  uint64
	LoginAt time.Time
}

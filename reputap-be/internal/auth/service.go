package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kittu533/reputap-be/internal/config"
)

// 💡 Go Concept: Custom errors
// We define sentinel errors that the handler can check with errors.Is().
// This separates "what went wrong" (service) from "what HTTP status to return" (handler).
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrPasswordTooWeak    = errors.New("password must be at least 8 characters")
	ErrWrongPassword      = errors.New("current password is incorrect")
)

// Service handles authentication business logic.
//
// 💡 Go Concept: Dependency Injection
// The service receives its dependencies (repo, config) via the constructor.
// This makes it testable — you can inject mock dependencies in tests.
type Service struct {
	repo *Repository
	cfg  *config.Config
}

// NewService creates a new auth service.
func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

// Login verifies credentials and returns JWT tokens.
//
// 💡 Go Concept: Method on struct
// `(s *Service)` means this function belongs to the Service struct.
// You call it like: service.Login(ctx, req)
func (s *Service) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	// Step 1: Find user by email
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		// 💡 Go Concept: errors.Is()
		// sql.ErrNoRows is returned by SQLX when no row matches the query.
		// We wrap it into our own error so the handler doesn't need to know about SQL.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("finding user: %w", err)
	}

	// Step 2: Check if user is active
	if user.Status != "active" {
		return nil, ErrUserInactive
	}

	// Step 3: Verify password
	// 💡 bcrypt: CompareHashAndPassword checks if the plaintext password
	// matches the stored hash. It returns nil on success, error on mismatch.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Step 4: Generate tokens
	accessToken, err := s.generateToken(user, s.cfg.JWTAccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := s.generateToken(user, s.cfg.JWTRefreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	// Step 5: Update last login (fire-and-forget, don't block on error)
	go s.repo.UpdateLastLogin(context.Background(), user.ID)

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.JWTAccessExpiry.Seconds()),
		User: UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			RoleID:   user.RoleID,
			RoleName: user.RoleName,
			Status:   user.Status,
		},
	}, nil
}

// RefreshToken validates a refresh token and returns new tokens.
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Parse and validate the refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get fresh user data (in case role/status changed)
	user, err := s.repo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("finding user: %w", err)
	}

	if user.Status != "active" {
		return nil, ErrUserInactive
	}

	// Generate new tokens using internal user model
	userWithPass := &UserWithPassword{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		RoleID:   user.RoleID,
		RoleName: user.RoleName,
		Status:   user.Status,
	}

	newAccessToken, err := s.generateToken(userWithPass, s.cfg.JWTAccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	newRefreshToken, err := s.generateToken(userWithPass, s.cfg.JWTRefreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.cfg.JWTAccessExpiry.Seconds()),
		User:         *user,
	}, nil
}

// GetCurrentUser retrieves the current user's info.
func (s *Service) GetCurrentUser(ctx context.Context, userID uint64) (*UserInfo, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("finding user: %w", err)
	}
	return user, nil
}

// ChangePassword validates and updates a user's password.
func (s *Service) ChangePassword(ctx context.Context, userID uint64, req ChangePasswordRequest) error {
	if len(req.NewPassword) < 8 {
		return ErrPasswordTooWeak
	}

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("finding user: %w", err)
	}

	// Get the full user with password to verify current password
	userFull, err := s.repo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("finding user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFull.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return ErrWrongPassword
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	return s.repo.UpdatePassword(ctx, userID, string(hash))
}

// GetUserPermissions returns permission names for a user's role.
func (s *Service) GetUserPermissions(ctx context.Context, roleID uint64) ([]string, error) {
	return s.repo.GetUserPermissions(ctx, roleID)
}

// ValidateToken parses and validates a JWT token, returning the claims.
//
// 💡 Go Concept: JWT Validation
// jwt.Parse takes the token string and a key function.
// The key function returns the secret used to verify the signature.
// If the token is expired or tampered, Parse returns an error.
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC (what we used to sign)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	// Extract claims from the token
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// 💡 Go Concept: Type assertion
	// JWT claims are stored as map[string]interface{}.
	// We need to assert each value to its expected type.
	// The .(float64) is because JSON numbers are parsed as float64.
	return &Claims{
		UserID:   uint64(mapClaims["user_id"].(float64)),
		Email:    mapClaims["email"].(string),
		RoleID:   uint64(mapClaims["role_id"].(float64)),
		RoleName: mapClaims["role_name"].(string),
	}, nil
}

// generateToken creates a JWT token with user claims and expiry.
func (s *Service) generateToken(user *UserWithPassword, expiry time.Duration) (string, error) {
	// 💡 Go Concept: jwt.MapClaims
	// Claims are the payload of the JWT. We embed user info so we
	// don't need a database lookup on every authenticated request.
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"role_id":   user.RoleID,
		"role_name": user.RoleName,
		"exp":       time.Now().Add(expiry).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

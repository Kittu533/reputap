package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Repository handles all database operations for authentication.
//
// 💡 Go Concept: Repository pattern
// The repository is the ONLY layer that talks to the database.
// It uses SQLX which lets you write raw SQL but maps results to structs
// using `db:"column_name"` struct tags (see model.go).
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new auth repository.
//
// 💡 Go Concept: Constructor function
// Go doesn't have constructors. Instead, we use "New" functions
// that return a pointer to the struct. This is the standard pattern.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// FindUserByEmail retrieves a user by email, including their role name.
//
// 💡 Go Concept: Context
// context.Context is passed through the call chain for cancellation/timeouts.
// Every database operation should accept a context.
//
// 💡 Go Concept: Error handling
// Go doesn't have exceptions. Functions return (value, error).
// The caller MUST check `if err != nil` — this is idiomatic Go.
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	var user UserWithPassword

	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.role_id, r.name AS role_name, u.status
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.email = ?
	`

	// 💡 SQLX: GetContext maps a single row to a struct.
	// The `db:"column_name"` tags on UserWithPassword tell SQLX
	// which struct field maps to which SQL column.
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindUserByID retrieves user info by ID.
func (r *Repository) FindUserByID(ctx context.Context, id uint64) (*UserInfo, error) {
	var user UserInfo

	query := `
		SELECT u.id, u.name, u.email, u.role_id, r.name AS role_name, u.status
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = ?
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateLastLogin sets the last_login_at timestamp for a user.
func (r *Repository) UpdateLastLogin(ctx context.Context, userID uint64) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// GetUserPermissions retrieves all permission names for a given role.
func (r *Repository) GetUserPermissions(ctx context.Context, roleID uint64) ([]string, error) {
	var permissions []string

	query := `
		SELECT p.name
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
	`

	// 💡 SQLX: SelectContext maps multiple rows to a slice.
	// Here we're selecting into a slice of strings since we only
	// need the permission names.
	err := r.db.SelectContext(ctx, &permissions, query, roleID)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// UpdatePassword updates the password hash for a user.
func (r *Repository) UpdatePassword(ctx context.Context, userID uint64, passwordHash string) error {
	query := `UPDATE users SET password_hash = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	return err
}

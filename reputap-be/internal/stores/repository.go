package stores

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/kittu533/reputap-be/internal/database"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, store *Store) error {
	query := `
		INSERT INTO stores (customer_id, store_code, name, slug, category, address, google_business_name, review_url, review_url_verified, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	id, err := r.db.InsertID(ctx, query,
		store.CustomerID,
		store.StoreCode,
		store.Name,
		store.Slug,
		store.Category,
		store.Address,
		store.GoogleBusinessName,
		store.ReviewURL,
		store.ReviewURLVerified,
		store.Status,
	)
	if err != nil {
		return err
	}

	store.ID = id
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*Store, error) {
	var store Store

	query := `SELECT * FROM stores WHERE id = ?`

	err := r.db.GetContext(ctx, &store, query, id)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]Store, int, error) {
	var total int
	var stores []Store

	where := "WHERE 1=1"
	var args []interface{}

	if strings.TrimSpace(params.Search) != "" {
		where += " AND (s.name LIKE ? OR s.store_code LIKE ?)"
		search := "%" + params.Search + "%"
		args = append(args, search, search)
	}
	if params.CustomerID > 0 {
		where += " AND s.customer_id = ?"
		args = append(args, params.CustomerID)
	}
	if params.Status != "" {
		where += " AND s.status = ?"
		args = append(args, params.Status)
	}

	countQuery := "SELECT COUNT(*) FROM stores s " + where
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := "SELECT s.* FROM stores s " + where + " ORDER BY s.id DESC LIMIT ? OFFSET ?"
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &stores, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

func (r *Repository) Update(ctx context.Context, store *Store) error {
	query := `
		UPDATE stores
		SET name = ?, slug = ?, category = ?, address = ?, google_business_name = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		store.Name,
		store.Slug,
		store.Category,
		store.Address,
		store.GoogleBusinessName,
		store.ID,
	)
	return err
}

func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	query := `UPDATE stores SET status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *Repository) UpdateReviewURL(ctx context.Context, id uint64, reviewURL string, verified bool) error {
	query := `UPDATE stores SET review_url = ?, review_url_verified = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, reviewURL, verified, id)
	return err
}

func (r *Repository) GetNextCode(ctx context.Context) (string, error) {
	var maxCode *string

	query := `SELECT MAX(store_code) FROM stores`
	err := r.db.GetContext(ctx, &maxCode, query)
	if err != nil {
		return "", err
	}

	nextNum := 1
	if maxCode != nil {
		parts := strings.Split(*maxCode, "-")
		if len(parts) == 2 {
			num, err := strconv.Atoi(parts[1])
			if err == nil {
				nextNum = num + 1
			}
		}
	}

	return fmt.Sprintf("STORE-%03d", nextNum), nil
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	var orderCount int
	err := r.db.GetContext(ctx, &orderCount, `SELECT COUNT(*) FROM orders WHERE store_id = ?`, id)
	if err != nil {
		return err
	}
	if orderCount > 0 {
		return fmt.Errorf("cannot delete store with %d associated orders", orderCount)
	}

	// Unlink QR codes and product instances
	_, _ = r.db.ExecContext(ctx, `UPDATE qr_codes SET store_id = NULL, status = 'unassigned' WHERE store_id = ?`, id)
	_, _ = r.db.ExecContext(ctx, `UPDATE product_instances SET store_id = NULL WHERE store_id = ?`, id)

	_, err = r.db.ExecContext(ctx, `DELETE FROM stores WHERE id = ?`, id)
	return err
}

func (r *Repository) SlugExists(ctx context.Context, slug string, excludeID uint64) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM stores WHERE slug = ? AND id != ?`
	err := r.db.GetContext(ctx, &count, query, slug, excludeID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

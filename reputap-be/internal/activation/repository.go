package activation

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/kittu533/reputap-be/internal/database"
)

var (
	ErrNotFound      = errors.New("kode QR tidak ditemukan")
	ErrAlreadyActive = errors.New("kartu ini sudah diaktivasi sebelumnya")
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9-]+`)
var multipleHyphens = regexp.MustCompile(`-{2,}`)

func slugify(name string) string {
	slug := strings.ToLower(name)
	slug = nonAlphanumeric.ReplaceAllString(slug, "-")
	slug = multipleHyphens.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

type QRRecord struct {
	ID                uint64  `db:"id"`
	ProductInstanceID uint64  `db:"product_instance_id"`
	StoreID           *uint64 `db:"store_id"`
	Code              string  `db:"code"`
	Status            string  `db:"status"`
	StoreName         string  `db:"store_name"`
	ReviewURL         string  `db:"review_url"`
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*QRRecord, error) {
	query := `
		SELECT q.id, q.product_instance_id, q.store_id, q.code, q.status,
		       COALESCE(s.name, '') as store_name,
		       COALESCE(s.review_url, '') as review_url
		FROM qr_codes q
		LEFT JOIN stores s ON q.store_id = s.id
		WHERE q.code = ?
	`
	var rec QRRecord
	if err := r.db.GetContext(ctx, &rec, query, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (r *Repository) Activate(ctx context.Context, req ActivateRequest) (uint64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Check QR code within transaction
	var qrID, instanceID uint64
	var currentStatus string
	var existingStoreID sql.NullInt64

	checkQuery := `SELECT id, product_instance_id, store_id, status FROM qr_codes WHERE code = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, checkQuery, req.Code).Scan(&qrID, &instanceID, &existingStoreID, &currentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	if currentStatus == "active" && existingStoreID.Valid && existingStoreID.Int64 > 0 {
		return 0, ErrAlreadyActive
	}

	// 2. Find or create Customer
	var customerID uint64
	custFindQuery := `SELECT id FROM customers WHERE phone = $1 LIMIT 1`
	err = tx.QueryRowContext(ctx, custFindQuery, req.Phone).Scan(&customerID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		var count int
		tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers`).Scan(&count)
		custCode := fmt.Sprintf("CUST-%03d", count+1)

		insertCustQuery := `
			INSERT INTO customers (customer_code, name, business_name, email, phone, address, city)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`
		err = tx.QueryRowContext(ctx, insertCustQuery,
			custCode, req.OwnerName, req.BusinessName, req.Email, req.Phone, req.Address, req.City,
		).Scan(&customerID)
		if err != nil {
			return 0, fmt.Errorf("creating customer: %w", err)
		}
	}

	// 3. Create Store
	var storeCount int
	tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM stores`).Scan(&storeCount)
	storeCode := fmt.Sprintf("STORE-%03d", storeCount+1)
	baseSlug := slugify(req.BusinessName)
	if baseSlug == "" {
		baseSlug = "toko"
	}
	slug := fmt.Sprintf("%s-%d", baseSlug, storeCount+1)

	var storeID uint64
	insertStoreQuery := `
		INSERT INTO stores (customer_id, store_code, name, slug, category, address, google_business_name, review_url, review_url_verified, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true, 'active')
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, insertStoreQuery,
		customerID, storeCode, req.BusinessName, slug, req.Category, req.Address, req.BusinessName, req.ReviewURL,
	).Scan(&storeID)
	if err != nil {
		return 0, fmt.Errorf("creating store: %w", err)
	}

	// 4. Update QR Code
	updateQRQuery := `
		UPDATE qr_codes
		SET store_id = $1, status = 'active', updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	if _, err := tx.ExecContext(ctx, updateQRQuery, storeID, qrID); err != nil {
		return 0, fmt.Errorf("updating qr_code: %w", err)
	}

	// 5. Update Product Instance
	updateInstQuery := `
		UPDATE product_instances
		SET store_id = $1, nfc_status = 'configured', updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	if _, err := tx.ExecContext(ctx, updateInstQuery, storeID, instanceID); err != nil {
		return 0, fmt.Errorf("updating product_instance: %w", err)
	}

	// 6. Audit Log
	newValues, _ := json.Marshal(map[string]interface{}{
		"code":       req.Code,
		"store_id":   storeID,
		"store_name": req.BusinessName,
		"phone":      req.Phone,
	})
	tx.ExecContext(ctx, `
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, new_values)
		VALUES (1, 'CARD_ACTIVATED', 'qr_codes', $1, $2)
	`, qrID, string(newValues))

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("committing activation: %w", err)
	}

	return storeID, nil
}

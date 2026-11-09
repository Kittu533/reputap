package qrcodes

import (
	"context"
	"fmt"
	"strings"

	"github.com/kittu533/reputap-be/internal/database"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*QRCode, error) {
	var qr QRCode

	query := `SELECT * FROM qr_codes WHERE id = ?`

	err := r.db.GetContext(ctx, &qr, query, id)
	if err != nil {
		return nil, err
	}

	return &qr, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*QRCode, error) {
	var qr QRCode

	query := `SELECT * FROM qr_codes WHERE code = ?`

	err := r.db.GetContext(ctx, &qr, query, code)
	if err != nil {
		return nil, err
	}

	return &qr, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]QRCode, int, error) {
	var total int
	var qrCodes []QRCode

	where := "WHERE 1=1"
	var args []interface{}

	if params.StoreID > 0 {
		where += " AND q.store_id = ?"
		args = append(args, params.StoreID)
	}
	if params.Status != "" {
		where += " AND q.status = ?"
		args = append(args, params.Status)
	}
	if strings.TrimSpace(params.Search) != "" {
		where += " AND q.code LIKE ?"
		search := "%" + params.Search + "%"
		args = append(args, search)
	}

	countQuery := "SELECT COUNT(*) FROM qr_codes q " + where
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := "SELECT q.* FROM qr_codes q " + where + " ORDER BY q.id DESC LIMIT ? OFFSET ?"
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &qrCodes, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return qrCodes, total, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	query := `UPDATE qr_codes SET status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *Repository) Update(ctx context.Context, id uint64, storeID *uint64, status string) error {
	query := `UPDATE qr_codes SET store_id = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, storeID, status, id)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		UPDATE product_instances SET store_id = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = (SELECT product_instance_id FROM qr_codes WHERE id = ?)
	`, storeID, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	var instanceID uint64
	err := r.db.GetContext(ctx, &instanceID, `SELECT product_instance_id FROM qr_codes WHERE id = ?`, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM scan_logs WHERE qr_code_id = ?`, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM qr_codes WHERE id = ?`, id)
	if err != nil {
		return err
	}

	_, _ = tx.ExecContext(ctx, `DELETE FROM quality_controls WHERE product_instance_id = ?`, instanceID)
	_, _ = tx.ExecContext(ctx, `DELETE FROM production_jobs WHERE product_instance_id = ?`, instanceID)

	_, err = tx.ExecContext(ctx, `DELETE FROM product_instances WHERE id = ?`, instanceID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) BatchGenerate(ctx context.Context, count int, prefix string) ([]QRCode, error) {
	if count <= 0 {
		count = 20
	}
	if prefix == "" {
		prefix = "RP26"
	}

	var maxID int
	_ = r.db.GetContext(ctx, &maxID, `SELECT COALESCE(COUNT(*), 0) FROM product_instances WHERE serial_number LIKE ?`, prefix+"-%")

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var createdCodes []QRCode
	for i := 1; i <= count; i++ {
		num := maxID + i
		serial := fmt.Sprintf("%s-%04d", prefix, num)
		code := fmt.Sprintf("rpt-%04d", num)

		var instanceID uint64
		err := tx.QueryRowContext(ctx, `
			INSERT INTO product_instances (serial_number, nfc_status, production_status, qc_status)
			VALUES ($1, 'not_configured', 'completed', 'passed')
			ON CONFLICT (serial_number) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
			RETURNING id
		`, serial).Scan(&instanceID)
		if err != nil {
			return nil, err
		}

		var qr QRCode
		err = tx.QueryRowxContext(ctx, `
			INSERT INTO qr_codes (product_instance_id, code, status)
			VALUES ($1, $2, 'unassigned')
			ON CONFLICT (code) DO UPDATE SET status = EXCLUDED.status
			RETURNING id, product_instance_id, store_id, code, status, created_at, updated_at
		`, instanceID, code).StructScan(&qr)
		if err != nil {
			return nil, err
		}

		createdCodes = append(createdCodes, qr)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdCodes, nil
}

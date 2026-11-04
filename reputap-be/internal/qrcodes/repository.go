package qrcodes

import (
	"context"
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

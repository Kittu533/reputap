package redirect

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByCode(ctx context.Context, code string) (*QRWithStore, error) {
	query := `
		SELECT qr.id as qr_id, qr.code as qr_code, qr.status as qr_status,
		       s.id as store_id, s.name as store_name, s.status as store_status,
		       COALESCE(s.review_url, '') as review_url
		FROM qr_codes qr
		JOIN stores s ON qr.store_id = s.id
		WHERE qr.code = ?`

	var result QRWithStore
	if err := r.db.GetContext(ctx, &result, query, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *Repository) InsertScanLog(ctx context.Context, log *ScanLog) error {
	query := `
		INSERT INTO scan_logs (qr_code_id, source, device_type, user_agent, referrer)
		VALUES (:qr_code_id, :source, :device_type, :user_agent, :referrer)`

	_, err := r.db.NamedExecContext(ctx, query, log)
	return err
}

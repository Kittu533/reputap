package qc

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, qc *QualityControl) error {
	query := `
		INSERT INTO quality_controls (product_instance_id, inspected_by, result, checklist_data, evidence_paths, notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		qc.ProductInstanceID,
		qc.InspectedBy,
		qc.Result,
		qc.ChecklistData,
		qc.EvidencePaths,
		qc.Notes,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	qc.ID = uint64(id)
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*QCWithInstance, error) {
	var qc QCWithInstance

	query := `
		SELECT qc.*, pi.serial_number, pi.store_id
		FROM quality_controls qc
		JOIN product_instances pi ON qc.product_instance_id = pi.id
		WHERE qc.id = ?
	`

	err := r.db.GetContext(ctx, &qc, query, id)
	if err != nil {
		return nil, err
	}

	return &qc, nil
}

func (r *Repository) ListByInstanceID(ctx context.Context, instanceID uint64) ([]QualityControl, error) {
	var qcs []QualityControl

	query := `SELECT * FROM quality_controls WHERE product_instance_id = ? ORDER BY inspected_at DESC`

	err := r.db.SelectContext(ctx, &qcs, query, instanceID)
	if err != nil {
		return nil, err
	}

	return qcs, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]QCWithInstance, int, error) {
	var total int
	var qcs []QCWithInstance

	var conditions []string
	var args []interface{}

	if params.ProductInstanceID > 0 {
		conditions = append(conditions, "qc.product_instance_id = ?")
		args = append(args, params.ProductInstanceID)
	}

	if params.Result != "" {
		conditions = append(conditions, "qc.result = ?")
		args = append(args, params.Result)
	}

	baseWhere := ""
	if len(conditions) > 0 {
		baseWhere = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM quality_controls qc" + baseWhere
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := `
		SELECT qc.*, pi.serial_number, pi.store_id
		FROM quality_controls qc
		JOIN product_instances pi ON qc.product_instance_id = pi.id` + baseWhere + ` ORDER BY qc.id DESC LIMIT ? OFFSET ?`
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &qcs, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return qcs, total, nil
}

func (r *Repository) UpdateInstanceQCStatus(ctx context.Context, instanceID uint64, status string) error {
	query := `UPDATE product_instances SET qc_status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, instanceID)
	return err
}

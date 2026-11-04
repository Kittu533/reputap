package production

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, job *Job) error {
	query := `
		INSERT INTO production_jobs (product_instance_id, assigned_to, status, checklist_data, notes)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		job.ProductInstanceID,
		job.AssignedTo,
		job.Status,
		job.ChecklistData,
		job.Notes,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	job.ID = uint64(id)
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*JobWithInstance, error) {
	var job JobWithInstance

	query := `
		SELECT pj.*, pi.serial_number, pi.store_id
		FROM production_jobs pj
		JOIN product_instances pi ON pj.product_instance_id = pi.id
		WHERE pj.id = ?
	`

	err := r.db.GetContext(ctx, &job, query, id)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]JobWithInstance, int, error) {
	var total int
	var jobs []JobWithInstance

	var conditions []string
	var args []interface{}

	if params.Status != "" {
		conditions = append(conditions, "pj.status = ?")
		args = append(args, params.Status)
	}

	if params.AssignedTo > 0 {
		conditions = append(conditions, "pj.assigned_to = ?")
		args = append(args, params.AssignedTo)
	}

	if params.StoreID > 0 {
		conditions = append(conditions, "pi.store_id = ?")
		args = append(args, params.StoreID)
	}

	baseWhere := ""
	if len(conditions) > 0 {
		baseWhere = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM production_jobs pj JOIN product_instances pi ON pj.product_instance_id = pi.id" + baseWhere
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := "SELECT pj.*, pi.serial_number, pi.store_id FROM production_jobs pj JOIN product_instances pi ON pj.product_instance_id = pi.id" + baseWhere + " ORDER BY pj.id DESC LIMIT ? OFFSET ?"
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &jobs, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status string, notes *string, checklistData *json.RawMessage) error {
	query := `
		UPDATE production_jobs
		SET status = ?, notes = ?, checklist_data = ?,
			started_at = CASE WHEN ? = 'in_progress' AND started_at IS NULL THEN NOW() ELSE started_at END,
			completed_at = CASE WHEN ? IN ('completed', 'failed') THEN NOW() ELSE completed_at END
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, status, notes, checklistData, status, status, id)
	return err
}

func (r *Repository) UpdateProductInstanceStatus(ctx context.Context, instanceID uint64, status string) error {
	query := `UPDATE product_instances SET production_status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, instanceID)
	return err
}

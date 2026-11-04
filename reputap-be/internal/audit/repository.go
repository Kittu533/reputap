package audit

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

func (r *Repository) Create(ctx context.Context, log *AuditLog) error {
	query := `
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, old_values, new_values)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		log.UserID,
		log.Action,
		log.EntityType,
		log.EntityID,
		log.OldValues,
		log.NewValues,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	log.ID = uint64(id)
	return nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]AuditLog, int, error) {
	var total int
	var logs []AuditLog

	var conditions []string
	var args []interface{}

	if params.EntityType != "" {
		conditions = append(conditions, "entity_type = ?")
		args = append(args, params.EntityType)
	}

	if params.EntityID > 0 {
		conditions = append(conditions, "entity_id = ?")
		args = append(args, params.EntityID)
	}

	if params.UserID > 0 {
		conditions = append(conditions, "user_id = ?")
		args = append(args, params.UserID)
	}

	if params.Action != "" {
		conditions = append(conditions, "action = ?")
		args = append(args, params.Action)
	}

	if params.DateFrom != "" {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, params.DateFrom)
	}

	if params.DateTo != "" {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, params.DateTo)
	}

	baseWhere := ""
	if len(conditions) > 0 {
		baseWhere = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM audit_logs" + baseWhere
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := "SELECT * FROM audit_logs" + baseWhere + " ORDER BY id DESC LIMIT ? OFFSET ?"
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &logs, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

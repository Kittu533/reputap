package design

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

func (r *Repository) Create(ctx context.Context, asset *DesignAsset) error {
	query := `
		INSERT INTO design_assets (order_id, store_id, file_path, file_name, file_type, file_size, version, approval_status, uploaded_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	id, err := r.db.InsertID(ctx, query,
		asset.OrderID,
		asset.StoreID,
		asset.FilePath,
		asset.FileName,
		asset.FileType,
		asset.FileSize,
		asset.Version,
		asset.ApprovalStatus,
		asset.UploadedBy,
	)
	if err != nil {
		return err
	}

	asset.ID = id
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*DesignAsset, error) {
	var asset DesignAsset

	query := `SELECT * FROM design_assets WHERE id = ?`

	err := r.db.GetContext(ctx, &asset, query, id)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *Repository) ListByOrderID(ctx context.Context, orderID uint64) ([]DesignAsset, error) {
	var assets []DesignAsset

	query := `SELECT * FROM design_assets WHERE order_id = ? ORDER BY version DESC`

	err := r.db.SelectContext(ctx, &assets, query, orderID)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]DesignAsset, int, error) {
	var total int
	var assets []DesignAsset

	var conditions []string
	var args []interface{}

	if params.OrderID > 0 {
		conditions = append(conditions, "order_id = ?")
		args = append(args, params.OrderID)
	}

	if params.StoreID > 0 {
		conditions = append(conditions, "store_id = ?")
		args = append(args, params.StoreID)
	}

	if params.Status != "" {
		conditions = append(conditions, "approval_status = ?")
		args = append(args, params.Status)
	}

	baseWhere := ""
	if len(conditions) > 0 {
		baseWhere = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM design_assets" + baseWhere
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := "SELECT * FROM design_assets" + baseWhere + " ORDER BY id DESC LIMIT ? OFFSET ?"
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &assets, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

func (r *Repository) UpdateApprovalStatus(ctx context.Context, id uint64, status string) error {
	query := `UPDATE design_assets SET approval_status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *Repository) GetLatestVersionByOrderID(ctx context.Context, orderID uint64) (int, error) {
	var version *int

	query := `SELECT MAX(version) FROM design_assets WHERE order_id = ?`

	err := r.db.GetContext(ctx, &version, query, orderID)
	if err != nil {
		return 0, err
	}

	if version == nil {
		return 0, nil
	}

	return *version, nil
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM design_assets WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

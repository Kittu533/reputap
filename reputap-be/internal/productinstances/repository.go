package productinstances

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderItemInfo struct {
	ID        uint64 `db:"id"`
	OrderID   uint64 `db:"order_id"`
	ProductID uint64 `db:"product_id"`
	Quantity  int    `db:"quantity"`
	StoreID   uint64 `db:"store_id"`
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, instance *ProductInstance) error {
	query := `
		INSERT INTO product_instances (order_item_id, store_id, serial_number, nfc_status, production_status, qc_status)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(ctx, query,
		instance.OrderItemID,
		instance.StoreID,
		instance.SerialNumber,
		instance.NFCStatus,
		instance.ProductionStatus,
		instance.QCStatus,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	instance.ID = uint64(id)
	return nil
}

func (r *Repository) CreateQRCode(ctx context.Context, tx *sqlx.Tx, instanceID uint64, storeID uint64, code string) error {
	query := `
		INSERT INTO qr_codes (product_instance_id, store_id, code, status)
		VALUES (?, ?, ?, 'active')
	`

	_, err := tx.ExecContext(ctx, query, instanceID, storeID, code)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*InstanceWithQR, error) {
	var instance InstanceWithQR

	query := `
		SELECT pi.*, qr.id AS qr_code_id, qr.code AS qr_code, qr.status AS qr_status
		FROM product_instances pi
		LEFT JOIN qr_codes qr ON qr.product_instance_id = pi.id
		WHERE pi.id = ?
	`

	err := r.db.GetContext(ctx, &instance, query, id)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

func (r *Repository) ListByOrderID(ctx context.Context, orderID uint64) ([]InstanceWithQR, error) {
	var instances []InstanceWithQR

	query := `
		SELECT pi.*, qr.id AS qr_code_id, qr.code AS qr_code, qr.status AS qr_status
		FROM product_instances pi
		JOIN order_items oi ON pi.order_item_id = oi.id
		LEFT JOIN qr_codes qr ON qr.product_instance_id = pi.id
		WHERE oi.order_id = ?
		ORDER BY pi.id ASC
	`

	err := r.db.SelectContext(ctx, &instances, query, orderID)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]InstanceWithQR, int, error) {
	var total int
	var instances []InstanceWithQR

	var conditions []string
	var args []interface{}

	if params.StoreID > 0 {
		conditions = append(conditions, "pi.store_id = ?")
		args = append(args, params.StoreID)
	}

	if params.NFCStatus != "" {
		conditions = append(conditions, "pi.nfc_status = ?")
		args = append(args, params.NFCStatus)
	}

	if params.ProdStatus != "" {
		conditions = append(conditions, "pi.production_status = ?")
		args = append(args, params.ProdStatus)
	}

	baseWhere := ""
	if len(conditions) > 0 {
		baseWhere = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM product_instances pi" + baseWhere
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := `
		SELECT pi.*, qr.id AS qr_code_id, qr.code AS qr_code, qr.status AS qr_status
		FROM product_instances pi
		LEFT JOIN qr_codes qr ON qr.product_instance_id = pi.id` + baseWhere + ` ORDER BY pi.id DESC LIMIT ? OFFSET ?`
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &instances, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return instances, total, nil
}

func (r *Repository) UpdateNFCStatus(ctx context.Context, id uint64, status string, configuredBy *uint64) error {
	var query string
	var args []interface{}

	if status == "configured" || status == "tested" {
		query = `UPDATE product_instances SET nfc_status = ?, nfc_configured_at = NOW(), nfc_configured_by = ? WHERE id = ?`
		args = []interface{}{status, configuredBy, id}
	} else {
		query = `UPDATE product_instances SET nfc_status = ? WHERE id = ?`
		args = []interface{}{status, id}
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) GetOrderItemsByOrderID(ctx context.Context, orderID uint64) ([]OrderItemInfo, error) {
	var items []OrderItemInfo

	query := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, o.store_id
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE oi.order_id = ?
	`

	err := r.db.SelectContext(ctx, &items, query, orderID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) CountInstancesByOrderItemID(ctx context.Context, orderItemID uint64) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM product_instances WHERE order_item_id = ?`

	err := r.db.GetContext(ctx, &count, query, orderItemID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

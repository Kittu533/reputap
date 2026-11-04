package shipping

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

func (r *Repository) Create(ctx context.Context, shipment *Shipment) error {
	query := `
		INSERT INTO shipments (order_id, recipient_name, recipient_address, recipient_phone, courier, tracking_number)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	id, err := r.db.InsertID(ctx, query,
		shipment.OrderID,
		shipment.RecipientName,
		shipment.RecipientAddress,
		shipment.RecipientPhone,
		shipment.Courier,
		shipment.TrackingNumber,
	)
	if err != nil {
		return err
	}

	shipment.ID = id
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*Shipment, error) {
	var shipment Shipment

	query := `SELECT * FROM shipments WHERE id = ?`

	err := r.db.GetContext(ctx, &shipment, query, id)
	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *Repository) GetByOrderID(ctx context.Context, orderID uint64) (*Shipment, error) {
	var shipment Shipment

	query := `SELECT * FROM shipments WHERE order_id = ?`

	err := r.db.GetContext(ctx, &shipment, query, orderID)
	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]Shipment, int, error) {
	var total int
	var shipments []Shipment

	var conditions []string
	var args []interface{}

	if params.OrderID > 0 {
		conditions = append(conditions, "order_id = ?")
		args = append(args, params.OrderID)
	}

	if params.Status != "" {
		conditions = append(conditions, "shipping_status = ?")
		args = append(args, params.Status)
	}

	baseWhere := ""
	if len(conditions) > 0 {
		baseWhere = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM shipments" + baseWhere
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := "SELECT * FROM shipments" + baseWhere + " ORDER BY id DESC LIMIT ? OFFSET ?"
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &shipments, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status string, trackingNumber *string, courier *string) error {
	var setClauses []string
	var args []interface{}

	setClauses = append(setClauses, "shipping_status = ?")
	args = append(args, status)

	if trackingNumber != nil {
		setClauses = append(setClauses, "tracking_number = ?")
		args = append(args, *trackingNumber)
	}

	if courier != nil {
		setClauses = append(setClauses, "courier = ?")
		args = append(args, *courier)
	}

	if status == "shipped" {
		setClauses = append(setClauses, "shipped_at = NOW()")
	}

	if status == "delivered" {
		setClauses = append(setClauses, "delivered_at = NOW()")
	}

	query := "UPDATE shipments SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

package orders

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, order *Order) error {
	query := `
		INSERT INTO orders (order_number, customer_id, store_id, order_status, payment_status, subtotal, discount, total, notes, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		order.OrderNumber,
		order.CustomerID,
		order.StoreID,
		order.OrderStatus,
		order.PaymentStatus,
		order.Subtotal,
		order.Discount,
		order.Total,
		order.Notes,
		order.CreatedBy,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	order.ID = uint64(id)
	return nil
}

func (r *Repository) CreateItems(ctx context.Context, items []OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
		VALUES (?, ?, ?, ?, ?)
	`

	for i := range items {
		result, err := r.db.ExecContext(ctx, query,
			items[i].OrderID,
			items[i].ProductID,
			items[i].Quantity,
			items[i].UnitPrice,
			items[i].Subtotal,
		)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		items[i].ID = uint64(id)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*Order, error) {
	var order Order

	query := `SELECT * FROM orders WHERE id = ?`

	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *Repository) GetItemsByOrderID(ctx context.Context, orderID uint64) ([]OrderItem, error) {
	var items []OrderItem

	query := `SELECT * FROM order_items WHERE order_id = ?`

	err := r.db.SelectContext(ctx, &items, query, orderID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]Order, int, error) {
	var total int
	var orders []Order

	var conditions []string
	var args []interface{}

	if strings.TrimSpace(params.Search) != "" {
		pattern := "%" + params.Search + "%"
		conditions = append(conditions, "order_number LIKE ?")
		args = append(args, pattern)
	}

	if params.CustomerID > 0 {
		conditions = append(conditions, "customer_id = ?")
		args = append(args, params.CustomerID)
	}

	if params.StoreID > 0 {
		conditions = append(conditions, "store_id = ?")
		args = append(args, params.StoreID)
	}

	if params.OrderStatus != "" {
		conditions = append(conditions, "order_status = ?")
		args = append(args, params.OrderStatus)
	}

	if params.PaymentStatus != "" {
		conditions = append(conditions, "payment_status = ?")
		args = append(args, params.PaymentStatus)
	}

	baseWhere := ""
	if len(conditions) > 0 {
		baseWhere = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM orders" + baseWhere
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := "SELECT * FROM orders" + baseWhere + " ORDER BY id DESC LIMIT ? OFFSET ?"
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &orders, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	query := `UPDATE orders SET order_status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *Repository) UpdatePaymentStatus(ctx context.Context, id uint64, status string) error {
	query := `UPDATE orders SET payment_status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *Repository) UpdateTotals(ctx context.Context, id uint64, subtotal, discount, total float64) error {
	query := `UPDATE orders SET subtotal = ?, discount = ?, total = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, subtotal, discount, total, id)
	return err
}

func (r *Repository) CreateStatusHistory(ctx context.Context, history *StatusHistory) error {
	query := `
		INSERT INTO order_status_history (order_id, from_status, to_status, changed_by, notes)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		history.OrderID,
		history.FromStatus,
		history.ToStatus,
		history.ChangedBy,
		history.Notes,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	history.ID = uint64(id)
	return nil
}

func (r *Repository) GetStatusHistory(ctx context.Context, orderID uint64) ([]StatusHistory, error) {
	var history []StatusHistory

	query := `SELECT * FROM order_status_history WHERE order_id = ? ORDER BY created_at ASC`

	err := r.db.SelectContext(ctx, &history, query, orderID)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (r *Repository) GenerateOrderNumber(ctx context.Context) (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("ORD-%s-", now.Format("200601"))

	var maxNumber *string

	query := `SELECT MAX(order_number) FROM orders WHERE order_number LIKE ?`
	err := r.db.GetContext(ctx, &maxNumber, query, prefix+"%")
	if err != nil {
		return "", err
	}

	nextNum := 1
	if maxNumber != nil {
		parts := strings.Split(*maxNumber, "-")
		if len(parts) == 3 {
			num, err := strconv.Atoi(parts[2])
			if err == nil {
				nextNum = num + 1
			}
		}
	}

	return fmt.Sprintf("%s%04d", prefix, nextNum), nil
}

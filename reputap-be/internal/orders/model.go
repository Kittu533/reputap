package orders

import "time"

type Order struct {
	ID            uint64    `json:"id" db:"id"`
	OrderNumber   string    `json:"order_number" db:"order_number"`
	CustomerID    uint64    `json:"customer_id" db:"customer_id"`
	StoreID       uint64    `json:"store_id" db:"store_id"`
	OrderStatus   string    `json:"order_status" db:"order_status"`
	PaymentStatus string    `json:"payment_status" db:"payment_status"`
	Subtotal      float64   `json:"subtotal" db:"subtotal"`
	Discount      float64   `json:"discount" db:"discount"`
	Total         float64   `json:"total" db:"total"`
	Notes         *string   `json:"notes" db:"notes"`
	CreatedBy     uint64    `json:"created_by" db:"created_by"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type OrderDetail struct {
	Order
	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        uint64  `json:"id" db:"id"`
	OrderID   uint64  `json:"order_id" db:"order_id"`
	ProductID uint64  `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
	Subtotal  float64 `json:"subtotal" db:"subtotal"`
}

type StatusHistory struct {
	ID         uint64    `json:"id" db:"id"`
	OrderID    uint64    `json:"order_id" db:"order_id"`
	FromStatus *string   `json:"from_status" db:"from_status"`
	ToStatus   string    `json:"to_status" db:"to_status"`
	ChangedBy  uint64    `json:"changed_by" db:"changed_by"`
	Notes      *string   `json:"notes" db:"notes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type CreateItemRequest struct {
	ProductID uint64  `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	UnitPrice float64 `json:"unit_price" validate:"required,gt=0"`
}

type CreateRequest struct {
	CustomerID uint64              `json:"customer_id" validate:"required"`
	StoreID    uint64              `json:"store_id" validate:"required"`
	Items      []CreateItemRequest `json:"items" validate:"required,min=1,dive"`
	Discount   float64             `json:"discount" validate:"gte=0"`
	Notes      *string             `json:"notes"`
}

type UpdateStatusRequest struct {
	Status string  `json:"status" validate:"required,oneof=pending confirmed processing production qc ready_to_ship shipped completed cancelled on_hold"`
	Notes  *string `json:"notes"`
}

type UpdatePaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status" validate:"required,oneof=unpaid partial paid refunded"`
}

type ListParams struct {
	Page          int
	PerPage       int
	Offset        int
	Search        string
	CustomerID    uint64
	StoreID       uint64
	OrderStatus   string
	PaymentStatus string
}

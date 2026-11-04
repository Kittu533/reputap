package shipping

import "time"

type Shipment struct {
	ID               uint64     `json:"id" db:"id"`
	OrderID          uint64     `json:"order_id" db:"order_id"`
	RecipientName    string     `json:"recipient_name" db:"recipient_name"`
	RecipientAddress string     `json:"recipient_address" db:"recipient_address"`
	RecipientPhone   *string    `json:"recipient_phone" db:"recipient_phone"`
	Courier          *string    `json:"courier" db:"courier"`
	TrackingNumber   *string    `json:"tracking_number" db:"tracking_number"`
	ShippingStatus   string     `json:"shipping_status" db:"shipping_status"`
	ShippedAt        *time.Time `json:"shipped_at" db:"shipped_at"`
	DeliveredAt      *time.Time `json:"delivered_at" db:"delivered_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateRequest struct {
	OrderID          uint64  `json:"order_id" validate:"required"`
	RecipientName    string  `json:"recipient_name" validate:"required,max=150"`
	RecipientAddress string  `json:"recipient_address" validate:"required"`
	RecipientPhone   *string `json:"recipient_phone" validate:"omitempty,max=20"`
	Courier          *string `json:"courier" validate:"omitempty,max=100"`
	TrackingNumber   *string `json:"tracking_number" validate:"omitempty,max=100"`
}

type UpdateStatusRequest struct {
	Status         string  `json:"status" validate:"required,oneof=waiting packed shipped delivered failed"`
	TrackingNumber *string `json:"tracking_number" validate:"omitempty,max=100"`
	Courier        *string `json:"courier" validate:"omitempty,max=100"`
}

type ListParams struct {
	Page    int
	PerPage int
	Offset  int
	OrderID uint64
	Status  string
}

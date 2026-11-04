package qrcodes

import "time"

type QRCode struct {
	ID                uint64    `json:"id" db:"id"`
	ProductInstanceID uint64    `json:"product_instance_id" db:"product_instance_id"`
	StoreID           uint64    `json:"store_id" db:"store_id"`
	Code              string    `json:"code" db:"code"`
	Status            string    `json:"status" db:"status"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// QRCodeDetail includes the computed redirect URL
type QRCodeDetail struct {
	QRCode
	RedirectURL string `json:"redirect_url"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended retired"`
}

type ListParams struct {
	Page    int
	PerPage int
	Offset  int
	StoreID uint64
	Status  string
	Search  string
}

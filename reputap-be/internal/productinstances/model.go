package productinstances

import "time"

type ProductInstance struct {
	ID               uint64     `json:"id" db:"id"`
	OrderItemID      uint64     `json:"order_item_id" db:"order_item_id"`
	StoreID          uint64     `json:"store_id" db:"store_id"`
	SerialNumber     string     `json:"serial_number" db:"serial_number"`
	NFCStatus        string     `json:"nfc_status" db:"nfc_status"`
	NFCConfiguredAt  *time.Time `json:"nfc_configured_at" db:"nfc_configured_at"`
	NFCConfiguredBy  *uint64    `json:"nfc_configured_by" db:"nfc_configured_by"`
	ProductionStatus string     `json:"production_status" db:"production_status"`
	QCStatus         string     `json:"qc_status" db:"qc_status"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type InstanceWithQR struct {
	ProductInstance
	QRCodeID uint64 `json:"qr_code_id" db:"qr_code_id"`
	QRCode   string `json:"qr_code" db:"qr_code"`
	QRStatus string `json:"qr_status" db:"qr_status"`
}

type UpdateNFCStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=not_configured configured tested failed replaced"`
}

type ListParams struct {
	Page       int
	PerPage    int
	Offset     int
	OrderID    uint64
	StoreID    uint64
	NFCStatus  string
	ProdStatus string
}

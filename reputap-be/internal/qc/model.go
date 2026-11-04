package qc

import (
	"encoding/json"
	"time"
)

type QualityControl struct {
	ID                uint64           `json:"id" db:"id"`
	ProductInstanceID uint64           `json:"product_instance_id" db:"product_instance_id"`
	InspectedBy       uint64           `json:"inspected_by" db:"inspected_by"`
	Result            string           `json:"result" db:"result"`
	ChecklistData     json.RawMessage  `json:"checklist_data" db:"checklist_data" swaggertype:"object"`
	EvidencePaths     *json.RawMessage `json:"evidence_paths" db:"evidence_paths" swaggertype:"object"`
	Notes             *string          `json:"notes" db:"notes"`
	InspectedAt       time.Time        `json:"inspected_at" db:"inspected_at"`
}

type QCWithInstance struct {
	QualityControl
	SerialNumber string `json:"serial_number" db:"serial_number"`
	StoreID      uint64 `json:"store_id" db:"store_id"`
}

type CreateRequest struct {
	ProductInstanceID uint64           `json:"product_instance_id" validate:"required"`
	Result            string           `json:"result" validate:"required,oneof=passed failed rework_required"`
	ChecklistData     json.RawMessage  `json:"checklist_data" validate:"required" swaggertype:"object"`
	EvidencePaths     *json.RawMessage `json:"evidence_paths" swaggertype:"object"`
	Notes             *string          `json:"notes"`
}

type ListParams struct {
	Page              int
	PerPage           int
	Offset            int
	ProductInstanceID uint64
	Result            string
}

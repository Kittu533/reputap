package production

import (
	"encoding/json"
	"time"
)

type Job struct {
	ID                uint64           `json:"id" db:"id"`
	ProductInstanceID uint64           `json:"product_instance_id" db:"product_instance_id"`
	AssignedTo        *uint64          `json:"assigned_to" db:"assigned_to"`
	Status            string           `json:"status" db:"status"`
	ChecklistData     *json.RawMessage `json:"checklist_data" db:"checklist_data" swaggertype:"object"`
	StartedAt         *time.Time       `json:"started_at" db:"started_at"`
	CompletedAt       *time.Time       `json:"completed_at" db:"completed_at"`
	Notes             *string          `json:"notes" db:"notes"`
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at" db:"updated_at"`
}

type JobWithInstance struct {
	Job
	SerialNumber string `json:"serial_number" db:"serial_number"`
	StoreID      uint64 `json:"store_id" db:"store_id"`
}

type CreateRequest struct {
	ProductInstanceID uint64           `json:"product_instance_id" validate:"required"`
	AssignedTo        *uint64          `json:"assigned_to"`
	Notes             *string          `json:"notes"`
	ChecklistData     *json.RawMessage `json:"checklist_data" swaggertype:"object"`
}

type UpdateStatusRequest struct {
	Status        string           `json:"status" validate:"required,oneof=pending in_progress completed failed"`
	Notes         *string          `json:"notes"`
	ChecklistData *json.RawMessage `json:"checklist_data" swaggertype:"object"`
}

type ListParams struct {
	Page       int
	PerPage    int
	Offset     int
	Status     string
	AssignedTo uint64
	StoreID    uint64
}

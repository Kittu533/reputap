package audit

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID         uint64           `json:"id" db:"id"`
	UserID     uint64           `json:"user_id" db:"user_id"`
	Action     string           `json:"action" db:"action"`
	EntityType string           `json:"entity_type" db:"entity_type"`
	EntityID   *uint64          `json:"entity_id" db:"entity_id"`
	OldValues  *json.RawMessage `json:"old_values" db:"old_values" swaggertype:"object"`
	NewValues  *json.RawMessage `json:"new_values" db:"new_values" swaggertype:"object"`
	CreatedAt  time.Time        `json:"created_at" db:"created_at"`
}

type ListParams struct {
	Page       int
	PerPage    int
	Offset     int
	EntityType string
	EntityID   uint64
	UserID     uint64
	Action     string
	DateFrom   string
	DateTo     string
}

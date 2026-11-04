package design

import "time"

type DesignAsset struct {
	ID             uint64    `json:"id" db:"id"`
	OrderID        uint64    `json:"order_id" db:"order_id"`
	StoreID        uint64    `json:"store_id" db:"store_id"`
	FilePath       string    `json:"file_path" db:"file_path"`
	FileName       string    `json:"file_name" db:"file_name"`
	FileType       string    `json:"file_type" db:"file_type"`
	FileSize       uint64    `json:"file_size" db:"file_size"`
	Version        uint      `json:"version" db:"version"`
	ApprovalStatus string    `json:"approval_status" db:"approval_status"`
	UploadedBy     uint64    `json:"uploaded_by" db:"uploaded_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type UpdateApprovalRequest struct {
	Status string `json:"status" validate:"required,oneof=draft pending_approval approved revision_required"`
}

type ListParams struct {
	Page    int
	PerPage int
	Offset  int
	OrderID uint64
	StoreID uint64
	Status  string
}

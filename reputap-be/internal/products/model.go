package products

import "time"

type Product struct {
	ID             uint64    `json:"id" db:"id"`
	SKU            string    `json:"sku" db:"sku"`
	Name           string    `json:"name" db:"name"`
	Description    *string   `json:"description" db:"description"`
	Material       *string   `json:"material" db:"material"`
	WidthMM        *float64  `json:"width_mm" db:"width_mm"`
	HeightMM       *float64  `json:"height_mm" db:"height_mm"`
	ThicknessMM    *float64  `json:"thickness_mm" db:"thickness_mm"`
	NFCChipType    *string   `json:"nfc_chip_type" db:"nfc_chip_type"`
	Finishing      *string   `json:"finishing" db:"finishing"`
	ProductionCost *float64  `json:"production_cost" db:"production_cost"`
	SellingPrice   *float64  `json:"selling_price" db:"selling_price"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRequest struct {
	SKU            string   `json:"sku" validate:"required,max=50"`
	Name           string   `json:"name" validate:"required,max=200"`
	Description    *string  `json:"description"`
	Material       *string  `json:"material" validate:"omitempty,max=100"`
	WidthMM        *float64 `json:"width_mm" validate:"omitempty,gte=0"`
	HeightMM       *float64 `json:"height_mm" validate:"omitempty,gte=0"`
	ThicknessMM    *float64 `json:"thickness_mm" validate:"omitempty,gte=0"`
	NFCChipType    *string  `json:"nfc_chip_type" validate:"omitempty,max=50"`
	Finishing      *string  `json:"finishing" validate:"omitempty,max=50"`
	ProductionCost *float64 `json:"production_cost" validate:"omitempty,gte=0"`
	SellingPrice   *float64 `json:"selling_price" validate:"omitempty,gte=0"`
	Status         *string  `json:"status" validate:"omitempty,oneof=active inactive discontinued"`
}

type UpdateRequest struct {
	Name           string   `json:"name" validate:"required,max=200"`
	SKU            *string  `json:"sku" validate:"omitempty,max=50"`
	Description    *string  `json:"description"`
	Material       *string  `json:"material" validate:"omitempty,max=100"`
	WidthMM        *float64 `json:"width_mm" validate:"omitempty,gte=0"`
	HeightMM       *float64 `json:"height_mm" validate:"omitempty,gte=0"`
	ThicknessMM    *float64 `json:"thickness_mm" validate:"omitempty,gte=0"`
	NFCChipType    *string  `json:"nfc_chip_type" validate:"omitempty,max=50"`
	Finishing      *string  `json:"finishing" validate:"omitempty,max=50"`
	ProductionCost *float64 `json:"production_cost" validate:"omitempty,gte=0"`
	SellingPrice   *float64 `json:"selling_price" validate:"omitempty,gte=0"`
	Status         *string  `json:"status" validate:"omitempty,oneof=active inactive discontinued"`
}

type ListParams struct {
	Page    int
	PerPage int
	Offset  int
	Search  string
	Status  string
}

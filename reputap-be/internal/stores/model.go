package stores

import "time"

type Store struct {
	ID                 uint64    `json:"id" db:"id"`
	CustomerID         uint64    `json:"customer_id" db:"customer_id"`
	StoreCode          string    `json:"store_code" db:"store_code"`
	Name               string    `json:"name" db:"name"`
	Slug               string    `json:"slug" db:"slug"`
	Category           string    `json:"category" db:"category"`
	Address            *string   `json:"address" db:"address"`
	GoogleBusinessName *string   `json:"google_business_name" db:"google_business_name"`
	ReviewURL          *string   `json:"review_url" db:"review_url"`
	ReviewURLVerified  bool      `json:"review_url_verified" db:"review_url_verified"`
	Status             string    `json:"status" db:"status"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRequest struct {
	CustomerID         uint64  `json:"customer_id" validate:"required"`
	Name               string  `json:"name" validate:"required,max=200"`
	Category           string  `json:"category" validate:"omitempty,oneof=coffee_shop restaurant barbershop salon clinic retail hotel other"`
	Address            *string `json:"address"`
	GoogleBusinessName *string `json:"google_business_name" validate:"omitempty,max=255"`
	ReviewURL          *string `json:"review_url" validate:"omitempty,max=500"`
}

type UpdateRequest struct {
	Name               string  `json:"name" validate:"required,max=200"`
	Category           string  `json:"category" validate:"omitempty,oneof=coffee_shop restaurant barbershop salon clinic retail hotel other"`
	Address            *string `json:"address"`
	GoogleBusinessName *string `json:"google_business_name" validate:"omitempty,max=255"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended"`
}

type UpdateReviewURLRequest struct {
	ReviewURL string `json:"review_url" validate:"required,max=500"`
}

type ListParams struct {
	Page       int
	PerPage    int
	Offset     int
	Search     string
	CustomerID uint64
	Status     string
}

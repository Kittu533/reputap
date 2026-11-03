package customers

import "time"

type Customer struct {
	ID           uint64    `json:"id" db:"id"`
	CustomerCode string    `json:"customer_code" db:"customer_code"`
	Name         string    `json:"name" db:"name"`
	BusinessName *string   `json:"business_name" db:"business_name"`
	Email        *string   `json:"email" db:"email"`
	Phone        *string   `json:"phone" db:"phone"`
	Address      *string   `json:"address" db:"address"`
	City         *string   `json:"city" db:"city"`
	Province     *string   `json:"province" db:"province"`
	Notes        *string   `json:"notes" db:"notes"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRequest struct {
	Name         string  `json:"name" validate:"required,max=150"`
	BusinessName *string `json:"business_name" validate:"omitempty,max=200"`
	Email        *string `json:"email" validate:"omitempty,email,max=255"`
	Phone        *string `json:"phone" validate:"omitempty,max=20"`
	Address      *string `json:"address"`
	City         *string `json:"city" validate:"omitempty,max=100"`
	Province     *string `json:"province" validate:"omitempty,max=100"`
	Notes        *string `json:"notes"`
}

type UpdateRequest struct {
	Name         string  `json:"name" validate:"required,max=150"`
	BusinessName *string `json:"business_name" validate:"omitempty,max=200"`
	Email        *string `json:"email" validate:"omitempty,email,max=255"`
	Phone        *string `json:"phone" validate:"omitempty,max=20"`
	Address      *string `json:"address"`
	City         *string `json:"city" validate:"omitempty,max=100"`
	Province     *string `json:"province" validate:"omitempty,max=100"`
	Notes        *string `json:"notes"`
}

type ListParams struct {
	Page    int
	PerPage int
	Offset  int
	Search  string
}

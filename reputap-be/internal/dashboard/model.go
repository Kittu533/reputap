package dashboard

import "time"

type OrderSummary struct {
	Total        int            `json:"total"`
	TotalRevenue float64        `json:"total_revenue"`
	ByStatus     map[string]int `json:"by_status"`
}

type CustomerSummary struct {
	Total int `json:"total"`
}

type StoreSummary struct {
	Total             int            `json:"total"`
	Active            int            `json:"active"`
	ReviewURLVerified int            `json:"review_url_verified"`
	ByStatus          map[string]int `json:"by_status"`
}

type QRCodeSummary struct {
	Total    int            `json:"total"`
	Active   int            `json:"active"`
	ByStatus map[string]int `json:"by_status"`
}

type ScanSummary struct {
	Total        int            `json:"total"`
	QRCount      int            `json:"qr_count"`
	NFCCount     int            `json:"nfc_count"`
	UnknownCount int            `json:"unknown_count"`
	ByDevice     map[string]int `json:"by_device"`
}

type SummaryResponse struct {
	Orders              OrderSummary    `json:"orders"`
	Customers           CustomerSummary `json:"customers"`
	Stores              StoreSummary    `json:"stores"`
	QRCodes             QRCodeSummary   `json:"qr_codes"`
	Scans               ScanSummary     `json:"scans"`
	RedirectSuccessRate float64         `json:"redirect_success_rate"`
}

type RecentActivity struct {
	ID         uint64    `json:"id" db:"id"`
	UserID     uint64    `json:"user_id" db:"user_id"`
	UserName   string    `json:"user_name" db:"user_name"`
	Action     string    `json:"action" db:"action"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   *uint64   `json:"entity_id" db:"entity_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

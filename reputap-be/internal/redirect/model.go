package redirect

import "time"

type QRWithStore struct {
	QRID        uint64 `db:"qr_id"`
	QRCode      string `db:"qr_code"`
	QRStatus    string `db:"qr_status"`
	StoreID     uint64 `db:"store_id"`
	StoreName   string `db:"store_name"`
	StoreStatus string `db:"store_status"`
	ReviewURL   string `db:"review_url"`
}

type ScanLog struct {
	ID         uint64    `db:"id"`
	QRCodeID   uint64    `db:"qr_code_id"`
	Source     string    `db:"source"`
	DeviceType string    `db:"device_type"`
	UserAgent  string    `db:"user_agent"`
	Referrer   string    `db:"referrer"`
	ScannedAt  time.Time `db:"scanned_at"`
}

type ResolveResult struct {
	Action    string // "redirect", "not_found", "inactive", "no_url"
	URL       string // the review URL to redirect to (only if redirect)
	QRCodeID  uint64
	StoreName string
}

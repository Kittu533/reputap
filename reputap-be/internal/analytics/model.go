package analytics

type StoreScanRank struct {
	StoreID      uint64 `json:"store_id" db:"store_id"`
	StoreName    string `json:"store_name" db:"store_name"`
	CustomerName string `json:"customer_name" db:"customer_name"`
	TotalScans   int    `json:"total_scans" db:"total_scans"`
	QRScans      int    `json:"qr_scans" db:"qr_scans"`
	NFCScans     int    `json:"nfc_scans" db:"nfc_scans"`
}

type OverviewResponse struct {
	TotalScans      int             `json:"total_scans"`
	QRScans         int             `json:"qr_scans"`
	NFCScans        int             `json:"nfc_scans"`
	UnknownScans    int             `json:"unknown_scans"`
	DeviceBreakdown map[string]int  `json:"device_breakdown"`
	TodayScans      int             `json:"today_scans"`
	YesterdayScans  int             `json:"yesterday_scans"`
	ThisWeekScans   int             `json:"this_week_scans"`
	ThisMonthScans  int             `json:"this_month_scans"`
	TopStores       []StoreScanRank `json:"top_stores"`
}

type TimeSeriesPoint struct {
	Date     string `json:"date" db:"date"`
	Total    int    `json:"total" db:"total"`
	QRCount  int    `json:"qr_count" db:"qr_count"`
	NFCCount int    `json:"nfc_count" db:"nfc_count"`
}

type ScansParams struct {
	Period  string // daily, weekly, monthly
	From    string // YYYY-MM-DD
	To      string // YYYY-MM-DD
	StoreID uint64
}

type StoreAnalyticsResponse struct {
	StoreID           uint64            `json:"store_id"`
	StoreName         string            `json:"store_name"`
	CustomerName      string            `json:"customer_name"`
	ReviewURL         string            `json:"review_url"`
	Status            string            `json:"status"`
	ReviewURLVerified bool              `json:"review_url_verified"`
	ActiveQRCodes     int               `json:"active_qr_codes"`
	TotalScans        int               `json:"total_scans"`
	QRScans           int               `json:"qr_scans"`
	NFCScans          int               `json:"nfc_scans"`
	DeviceBreakdown   map[string]int    `json:"device_breakdown"`
	DailyScans        []TimeSeriesPoint `json:"daily_scans"`
}

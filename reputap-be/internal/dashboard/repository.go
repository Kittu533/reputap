package dashboard

import (
	"context"

	"github.com/kittu533/reputap-be/internal/database"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetSummary(ctx context.Context) (*SummaryResponse, error) {
	summary := &SummaryResponse{
		Orders: OrderSummary{
			ByStatus: make(map[string]int),
		},
		Customers: CustomerSummary{},
		Stores: StoreSummary{
			ByStatus: make(map[string]int),
		},
		QRCodes: QRCodeSummary{
			ByStatus: make(map[string]int),
		},
		Scans: ScanSummary{
			ByDevice: make(map[string]int),
		},
		RedirectSuccessRate: 100.0,
	}

	// 1. Orders
	type orderStat struct {
		Status  string  `db:"order_status"`
		Count   int     `db:"cnt"`
		Revenue float64 `db:"rev"`
	}
	var orderStats []orderStat
	err := r.db.SelectContext(ctx, &orderStats, `
		SELECT order_status, COUNT(*) as cnt, COALESCE(SUM(total), 0) as rev
		FROM orders
		GROUP BY order_status
	`)
	if err != nil {
		return nil, err
	}
	for _, s := range orderStats {
		summary.Orders.Total += s.Count
		if s.Status != "cancelled" {
			summary.Orders.TotalRevenue += s.Revenue
		}
		summary.Orders.ByStatus[s.Status] = s.Count
	}

	// 2. Customers
	err = r.db.GetContext(ctx, &summary.Customers.Total, `SELECT COUNT(*) FROM customers`)
	if err != nil {
		return nil, err
	}

	// 3. Stores
	type storeStat struct {
		Status            string `db:"status"`
		ReviewURLVerified bool   `db:"review_url_verified"`
		Count             int    `db:"cnt"`
	}
	var storeStats []storeStat
	err = r.db.SelectContext(ctx, &storeStats, `
		SELECT status, review_url_verified, COUNT(*) as cnt
		FROM stores
		GROUP BY status, review_url_verified
	`)
	if err != nil {
		return nil, err
	}
	for _, s := range storeStats {
		summary.Stores.Total += s.Count
		summary.Stores.ByStatus[s.Status] += s.Count
		if s.Status == "active" {
			summary.Stores.Active += s.Count
		}
		if s.ReviewURLVerified {
			summary.Stores.ReviewURLVerified += s.Count
		}
	}

	// 4. QR Codes
	type qrStat struct {
		Status string `db:"status"`
		Count  int    `db:"cnt"`
	}
	var qrStats []qrStat
	err = r.db.SelectContext(ctx, &qrStats, `
		SELECT status, COUNT(*) as cnt
		FROM qr_codes
		GROUP BY status
	`)
	if err != nil {
		return nil, err
	}
	for _, s := range qrStats {
		summary.QRCodes.Total += s.Count
		summary.QRCodes.ByStatus[s.Status] = s.Count
		if s.Status == "active" {
			summary.QRCodes.Active += s.Count
		}
	}

	// 5. Scans
	type scanStat struct {
		Source     string `db:"source"`
		DeviceType string `db:"device_type"`
		Count      int    `db:"cnt"`
	}
	var scanStats []scanStat
	err = r.db.SelectContext(ctx, &scanStats, `
		SELECT source, device_type, COUNT(*) as cnt
		FROM scan_logs
		GROUP BY source, device_type
	`)
	if err != nil {
		return nil, err
	}
	for _, s := range scanStats {
		summary.Scans.Total += s.Count
		summary.Scans.ByDevice[s.DeviceType] += s.Count
		switch s.Source {
		case "qr":
			summary.Scans.QRCount += s.Count
		case "nfc":
			summary.Scans.NFCCount += s.Count
		default:
			summary.Scans.UnknownCount += s.Count
		}
	}

	if summary.QRCodes.Total > 0 {
		summary.RedirectSuccessRate = float64(summary.QRCodes.Active) / float64(summary.QRCodes.Total) * 100
	}

	return summary, nil
}

func (r *Repository) GetRecentActivities(ctx context.Context, limit int) ([]RecentActivity, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	var activities []RecentActivity
	query := `
		SELECT a.id, a.user_id, COALESCE(u.name, 'System') as user_name, a.action, a.entity_type, a.entity_id, a.created_at
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
		LIMIT ?
	`
	err := r.db.SelectContext(ctx, &activities, query, limit)
	if err != nil {
		return nil, err
	}
	if activities == nil {
		activities = []RecentActivity{}
	}
	return activities, nil
}

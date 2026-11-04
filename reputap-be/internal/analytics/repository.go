package analytics

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

var ErrStoreNotFound = errors.New("store not found")

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetOverview(ctx context.Context) (*OverviewResponse, error) {
	resp := &OverviewResponse{
		DeviceBreakdown: make(map[string]int),
		TopStores:       []StoreScanRank{},
	}

	// 1. Scan totals & device breakdown
	type scanStat struct {
		Source     string `db:"source"`
		DeviceType string `db:"device_type"`
		Count      int    `db:"cnt"`
	}
	var stats []scanStat
	err := r.db.SelectContext(ctx, &stats, `
		SELECT source, device_type, COUNT(*) as cnt
		FROM scan_logs
		GROUP BY source, device_type
	`)
	if err != nil {
		return nil, err
	}

	for _, s := range stats {
		resp.TotalScans += s.Count
		resp.DeviceBreakdown[s.DeviceType] += s.Count
		switch s.Source {
		case "qr":
			resp.QRScans += s.Count
		case "nfc":
			resp.NFCScans += s.Count
		default:
			resp.UnknownScans += s.Count
		}
	}

	// 2. Periods (today, yesterday, this week, this month)
	type periodStats struct {
		TodayScans     int `db:"today_scans"`
		YesterdayScans int `db:"yesterday_scans"`
		ThisWeekScans  int `db:"this_week_scans"`
		ThisMonthScans int `db:"this_month_scans"`
	}
	var pStats periodStats
	err = r.db.GetContext(ctx, &pStats, `
		SELECT
			COALESCE(SUM(CASE WHEN DATE(scanned_at) = CURDATE() THEN 1 ELSE 0 END), 0) as today_scans,
			COALESCE(SUM(CASE WHEN DATE(scanned_at) = CURDATE() - INTERVAL 1 DAY THEN 1 ELSE 0 END), 0) as yesterday_scans,
			COALESCE(SUM(CASE WHEN YEARWEEK(scanned_at, 1) = YEARWEEK(CURDATE(), 1) THEN 1 ELSE 0 END), 0) as this_week_scans,
			COALESCE(SUM(CASE WHEN YEAR(scanned_at) = YEAR(CURDATE()) AND MONTH(scanned_at) = MONTH(CURDATE()) THEN 1 ELSE 0 END), 0) as this_month_scans
		FROM scan_logs
	`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	resp.TodayScans = pStats.TodayScans
	resp.YesterdayScans = pStats.YesterdayScans
	resp.ThisWeekScans = pStats.ThisWeekScans
	resp.ThisMonthScans = pStats.ThisMonthScans

	// 3. Top 5 stores
	err = r.db.SelectContext(ctx, &resp.TopStores, `
		SELECT s.id as store_id, s.name as store_name, c.name as customer_name,
			COUNT(sl.id) as total_scans,
			COALESCE(SUM(CASE WHEN sl.source = 'qr' THEN 1 ELSE 0 END), 0) as qr_scans,
			COALESCE(SUM(CASE WHEN sl.source = 'nfc' THEN 1 ELSE 0 END), 0) as nfc_scans
		FROM scan_logs sl
		JOIN qr_codes q ON sl.qr_code_id = q.id
		JOIN stores s ON q.store_id = s.id
		JOIN customers c ON s.customer_id = c.id
		GROUP BY s.id, s.name, c.name
		ORDER BY total_scans DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	if resp.TopStores == nil {
		resp.TopStores = []StoreScanRank{}
	}

	return resp, nil
}

func (r *Repository) GetScansTimeSeries(ctx context.Context, params ScansParams) ([]TimeSeriesPoint, error) {
	var dateFormat string
	switch params.Period {
	case "weekly":
		dateFormat = "%x-W%v"
	case "monthly":
		dateFormat = "%Y-%m"
	default:
		dateFormat = "%Y-%m-%d"
	}

	var conditions []string
	var args []interface{}

	if params.StoreID > 0 {
		conditions = append(conditions, "q.store_id = ?")
		args = append(args, params.StoreID)
	}

	if params.From != "" {
		conditions = append(conditions, "sl.scanned_at >= ?")
		args = append(args, params.From+" 00:00:00")
	}

	if params.To != "" {
		conditions = append(conditions, "sl.scanned_at <= ?")
		args = append(args, params.To+" 23:59:59")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT
			DATE_FORMAT(sl.scanned_at, '%s') as date,
			COUNT(sl.id) as total,
			COALESCE(SUM(CASE WHEN sl.source = 'qr' THEN 1 ELSE 0 END), 0) as qr_count,
			COALESCE(SUM(CASE WHEN sl.source = 'nfc' THEN 1 ELSE 0 END), 0) as nfc_count
		FROM scan_logs sl
		JOIN qr_codes q ON sl.qr_code_id = q.id
		%s
		GROUP BY date
		ORDER BY date ASC
	`, dateFormat, whereClause)

	var points []TimeSeriesPoint
	err := r.db.SelectContext(ctx, &points, query, args...)
	if err != nil {
		return nil, err
	}
	if points == nil {
		points = []TimeSeriesPoint{}
	}
	return points, nil
}

func (r *Repository) GetStoreAnalytics(ctx context.Context, storeID uint64) (*StoreAnalyticsResponse, error) {
	resp := &StoreAnalyticsResponse{
		StoreID:         storeID,
		DeviceBreakdown: make(map[string]int),
		DailyScans:      []TimeSeriesPoint{},
	}

	// 1. Store details
	type storeRow struct {
		Name              string `db:"name"`
		CustomerName      string `db:"customer_name"`
		ReviewURL         string `db:"review_url"`
		Status            string `db:"status"`
		ReviewURLVerified bool   `db:"review_url_verified"`
	}
	var sr storeRow
	err := r.db.GetContext(ctx, &sr, `
		SELECT s.name, c.name as customer_name, s.review_url, s.status, s.review_url_verified
		FROM stores s
		JOIN customers c ON s.customer_id = c.id
		WHERE s.id = ?
	`, storeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStoreNotFound
		}
		return nil, err
	}
	resp.StoreName = sr.Name
	resp.CustomerName = sr.CustomerName
	resp.ReviewURL = sr.ReviewURL
	resp.Status = sr.Status
	resp.ReviewURLVerified = sr.ReviewURLVerified

	// 2. Active QR codes
	err = r.db.GetContext(ctx, &resp.ActiveQRCodes, `
		SELECT COUNT(*)
		FROM qr_codes
		WHERE store_id = ? AND status = 'active'
	`, storeID)
	if err != nil {
		return nil, err
	}

	// 3. Scan logs for this store
	type scanRow struct {
		Source     string `db:"source"`
		DeviceType string `db:"device_type"`
		Count      int    `db:"cnt"`
	}
	var scanRows []scanRow
	err = r.db.SelectContext(ctx, &scanRows, `
		SELECT sl.source, sl.device_type, COUNT(sl.id) as cnt
		FROM scan_logs sl
		JOIN qr_codes q ON sl.qr_code_id = q.id
		WHERE q.store_id = ?
		GROUP BY sl.source, sl.device_type
	`, storeID)
	if err != nil {
		return nil, err
	}
	for _, sr := range scanRows {
		resp.TotalScans += sr.Count
		resp.DeviceBreakdown[sr.DeviceType] += sr.Count
		switch sr.Source {
		case "qr":
			resp.QRScans += sr.Count
		case "nfc":
			resp.NFCScans += sr.Count
		}
	}

	// 4. Daily scans last 30 days
	err = r.db.SelectContext(ctx, &resp.DailyScans, `
		SELECT
			DATE_FORMAT(sl.scanned_at, '%Y-%m-%d') as date,
			COUNT(sl.id) as total,
			COALESCE(SUM(CASE WHEN sl.source = 'qr' THEN 1 ELSE 0 END), 0) as qr_count,
			COALESCE(SUM(CASE WHEN sl.source = 'nfc' THEN 1 ELSE 0 END), 0) as nfc_count
		FROM scan_logs sl
		JOIN qr_codes q ON sl.qr_code_id = q.id
		WHERE q.store_id = ? AND sl.scanned_at >= CURDATE() - INTERVAL 30 DAY
		GROUP BY date
		ORDER BY date ASC
	`, storeID)
	if err != nil {
		return nil, err
	}
	if resp.DailyScans == nil {
		resp.DailyScans = []TimeSeriesPoint{}
	}

	return resp, nil
}

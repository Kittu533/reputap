package redirect

import (
	"context"
	"log/slog"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Resolve(ctx context.Context, code string, userAgent string, referrer string, source string) (*ResolveResult, error) {
	qr, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if qr == nil {
		return &ResolveResult{Action: "not_found"}, nil
	}

	result := &ResolveResult{
		QRCodeID:  qr.QRID,
		StoreName: qr.StoreName,
	}

	if qr.QRStatus == "unassigned" || qr.StoreID == 0 {
		result.Action = "unassigned"
		return result, nil
	}

	if qr.QRStatus != "active" {
		result.Action = "inactive"
		return result, nil
	}

	if qr.StoreStatus != "active" {
		result.Action = "inactive"
		return result, nil
	}

	if qr.ReviewURL == "" {
		result.Action = "no_url"
		return result, nil
	}

	source = normalizeSource(source)
	deviceType := detectDevice(userAgent)

	result.Action = "redirect"
	result.URL = qr.ReviewURL

	go func() {
		scanLog := &ScanLog{
			QRCodeID:   result.QRCodeID,
			Source:     source,
			DeviceType: deviceType,
			UserAgent:  userAgent,
			Referrer:   referrer,
		}
		if err := s.repo.InsertScanLog(context.Background(), scanLog); err != nil {
			slog.Error("failed to log scan", "error", err, "qr_code_id", result.QRCodeID)
		}
	}()

	return result, nil
}

func normalizeSource(source string) string {
	switch source {
	case "qr", "nfc":
		return source
	default:
		return "unknown"
	}
}

func detectDevice(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "android"):
		return "android"
	case strings.Contains(ua, "iphone"), strings.Contains(ua, "ipad"):
		return "ios"
	default:
		return "other"
	}
}

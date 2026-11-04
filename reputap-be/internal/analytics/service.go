package analytics

import (
	"context"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOverview(ctx context.Context) (*OverviewResponse, error) {
	return s.repo.GetOverview(ctx)
}

func (s *Service) GetScansTimeSeries(ctx context.Context, params ScansParams) ([]TimeSeriesPoint, error) {
	if params.Period == "" {
		params.Period = "daily"
	}
	// Default to last 30 days if from/to are empty
	if params.From == "" && params.To == "" {
		now := time.Now()
		params.To = now.Format("2006-01-02")
		params.From = now.AddDate(0, 0, -30).Format("2006-01-02")
	}

	return s.repo.GetScansTimeSeries(ctx, params)
}

func (s *Service) GetStoreAnalytics(ctx context.Context, storeID uint64) (*StoreAnalyticsResponse, error) {
	return s.repo.GetStoreAnalytics(ctx, storeID)
}

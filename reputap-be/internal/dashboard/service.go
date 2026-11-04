package dashboard

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSummary(ctx context.Context) (*SummaryResponse, error) {
	return s.repo.GetSummary(ctx)
}

func (s *Service) GetRecentActivities(ctx context.Context, limit int) ([]RecentActivity, error) {
	return s.repo.GetRecentActivities(ctx, limit)
}
